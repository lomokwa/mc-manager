import './style.css'
import {
  ApiError,
  clearConfig,
  consoleWebSocketUrl,
  getApiBase,
  getApiKey,
  getStatus,
  listPlayers,
  setApiBase,
  setApiKey,
  startServer,
  stopServer,
  type Player,
} from './api.ts'

const app = document.querySelector<HTMLDivElement>('#app')!
const MAX_CONSOLE_LINES = 1000

interface DashRefs {
  conn: HTMLElement
  badge: HTMLElement
  startBtn: HTMLButtonElement
  stopBtn: HTMLButtonElement
  consoleEl: HTMLPreElement
  cmdInput: HTMLInputElement
  sendBtn: HTMLButtonElement
  playersBody: HTMLElement
}

let dash: DashRefs | undefined
let pollTimer: number | undefined
let ws: WebSocket | undefined

// ----------------------------------------------------------------- helpers

function h<K extends keyof HTMLElementTagNameMap>(
  tag: K,
  attrs: Record<string, string> = {},
  ...children: (Node | string)[]
): HTMLElementTagNameMap[K] {
  const node = document.createElement(tag)
  for (const [k, v] of Object.entries(attrs)) {
    if (k === 'class') node.className = v
    else node.setAttribute(k, v)
  }
  node.append(...children)
  return node
}

function field(labelText: string, input: HTMLElement): HTMLLabelElement {
  return h('label', { class: 'field' }, h('span', {}, labelText), input)
}

function toast(message: string, kind: 'error' | 'info' = 'info'): void {
  const t = h('div', { class: `toast toast-${kind}` }, message)
  document.body.append(t)
  requestAnimationFrame(() => t.classList.add('show'))
  setTimeout(() => {
    t.classList.remove('show')
    setTimeout(() => t.remove(), 300)
  }, 4500)
}

function handleError(err: unknown): void {
  if (err instanceof ApiError) {
    toast(err.message, 'error')
    if (err.status === 401) disconnect() // bad/expired key — back to connect
    return
  }
  toast((err as Error).message || 'unexpected error', 'error')
}

// ------------------------------------------------------------ connect view

function renderConnect(): void {
  stopPolling()

  const baseInput = h('input', {
    type: 'text',
    placeholder: 'blank = same origin (e.g. http://localhost:8080)',
  })
  baseInput.value = getApiBase()

  const keyInput = h('input', { type: 'password', placeholder: 'your API key' })
  keyInput.value = getApiKey()

  const connect = h('button', { class: 'primary' }, 'Connect')
  const submit = (): void => {
    const key = keyInput.value.trim()
    if (!key) {
      toast('An API key is required', 'error')
      return
    }
    setApiBase(baseInput.value.trim())
    setApiKey(key)
    renderDashboard()
  }
  connect.addEventListener('click', submit)
  keyInput.addEventListener('keydown', (e) => {
    if (e.key === 'Enter') submit()
  })

  app.replaceChildren(
    h(
      'div',
      { class: 'card connect' },
      h('h1', {}, 'MC Manager'),
      h('p', { class: 'muted' }, 'Connect to your mc-manager API to get started.'),
      field('API URL', baseInput),
      field('API key', keyInput),
      connect,
    ),
  )
  keyInput.focus()
}

// ---------------------------------------------------------- dashboard view

function renderDashboard(): void {
  const conn = h('span', { class: 'conn' }, 'idle')
  const changeKey = h('button', { class: 'link' }, 'Change key')
  changeKey.addEventListener('click', disconnect)
  const header = h(
    'header',
    {},
    h('h1', {}, 'MC Manager'),
    h('div', { class: 'header-right' }, conn, changeKey),
  )

  const badge = h('span', { class: 'badge badge-unknown' }, 'unknown')
  const startBtn = h('button', { class: 'primary' }, 'Start')
  const stopBtn = h('button', { class: 'danger' }, 'Stop')
  const refreshStatusBtn = h('button', { class: 'ghost' }, 'Refresh')
  startBtn.addEventListener('click', () => void doStart())
  stopBtn.addEventListener('click', () => void doStop())
  refreshStatusBtn.addEventListener('click', () => void refreshStatus())
  const statusCard = h(
    'section',
    { class: 'card' },
    h('h2', {}, 'Server'),
    h('div', { class: 'status-row' }, h('span', { class: 'muted' }, 'State:'), badge),
    h('div', { class: 'actions' }, startBtn, stopBtn, refreshStatusBtn),
  )

  const consoleEl = h('pre', { class: 'console' })
  const cmdInput = h('input', {
    class: 'cmd',
    placeholder: 'Send a command, e.g. say hello',
    autocomplete: 'off',
  })
  const sendBtn = h('button', { class: 'primary' }, 'Send')
  const cmdForm = h('form', { class: 'cmd-form' }, cmdInput, sendBtn)
  cmdForm.addEventListener('submit', (e) => {
    e.preventDefault()
    sendCommand()
  })
  const consoleCard = h(
    'section',
    { class: 'card console-card' },
    h('h2', {}, 'Console'),
    consoleEl,
    cmdForm,
  )

  const refreshPlayersBtn = h('button', { class: 'link' }, 'refresh')
  refreshPlayersBtn.addEventListener('click', () => void refreshPlayers())
  const playersBody = h('div', { class: 'players' })
  const playersCard = h(
    'section',
    { class: 'card' },
    h('h2', {}, 'Players ', refreshPlayersBtn),
    playersBody,
  )

  dash = { conn, badge, startBtn, stopBtn, consoleEl, cmdInput, sendBtn, playersBody }
  app.replaceChildren(header, h('main', {}, statusCard, consoleCard, playersCard))

  void refreshStatus()
  void refreshPlayers()
  startPolling()
}

// ------------------------------------------------------------ status + ops

function applyStatus(isRunning: boolean): void {
  if (!dash) return
  dash.badge.textContent = isRunning ? 'running' : 'stopped'
  dash.badge.className = `badge ${isRunning ? 'badge-up' : 'badge-down'}`
  dash.startBtn.disabled = isRunning
  dash.stopBtn.disabled = !isRunning
  dash.cmdInput.disabled = !isRunning
  dash.sendBtn.disabled = !isRunning
  if (isRunning) ensureConsole()
  else teardownConsole()
}

async function refreshStatus(): Promise<void> {
  try {
    applyStatus((await getStatus()).running)
  } catch (err) {
    handleError(err)
  }
}

async function doStart(): Promise<void> {
  if (!dash) return
  dash.startBtn.disabled = true
  appendConsole('> starting server…')
  try {
    await startServer()
    toast('Server starting', 'info')
  } catch (err) {
    handleError(err)
  } finally {
    await refreshStatus()
  }
}

async function doStop(): Promise<void> {
  if (!dash) return
  dash.stopBtn.disabled = true
  appendConsole('> stopping server…')
  try {
    await stopServer()
    toast('Server stopping', 'info')
  } catch (err) {
    handleError(err)
  } finally {
    await refreshStatus()
  }
}

// ---------------------------------------------------------------- console

function sendCommand(): void {
  if (!dash) return
  const cmd = dash.cmdInput.value.trim()
  if (!cmd) return
  if (ws && ws.readyState === WebSocket.OPEN) {
    ws.send(cmd)
    appendConsole(`> ${cmd}`)
    dash.cmdInput.value = ''
  } else {
    toast('Console is not connected', 'error')
  }
}

function ensureConsole(): void {
  if (ws && (ws.readyState === WebSocket.OPEN || ws.readyState === WebSocket.CONNECTING)) return
  const socket = new WebSocket(consoleWebSocketUrl())
  ws = socket
  setConn('connecting…')
  socket.addEventListener('open', () => setConn('live', true))
  socket.addEventListener('message', (e) => {
    if (typeof e.data === 'string') appendConsole(e.data)
  })
  socket.addEventListener('close', () => {
    if (ws === socket) {
      ws = undefined
      setConn('disconnected')
    }
  })
}

function teardownConsole(): void {
  setConn('idle')
  if (ws) {
    const socket = ws
    ws = undefined
    socket.close()
  }
}

function setConn(textLabel: string, live = false): void {
  if (!dash) return
  dash.conn.textContent = textLabel
  dash.conn.className = `conn ${live ? 'conn-live' : ''}`
}

function appendConsole(line: string): void {
  if (!dash) return
  const el = dash.consoleEl
  const atBottom = el.scrollHeight - el.scrollTop - el.clientHeight < 40
  el.append(line + '\n')
  while (el.childNodes.length > MAX_CONSOLE_LINES && el.firstChild) {
    el.removeChild(el.firstChild)
  }
  if (atBottom) el.scrollTop = el.scrollHeight
}

// ----------------------------------------------------------------- players

async function refreshPlayers(): Promise<void> {
  if (!dash) return
  try {
    renderPlayers(await listPlayers())
  } catch (err) {
    // Players are read from the server's usercache.json, which doesn't exist
    // until the server has run at least once.
    if (err instanceof ApiError && err.status === 500) {
      dash.playersBody.replaceChildren(h('p', { class: 'muted' }, 'No player data yet.'))
      return
    }
    handleError(err)
  }
}

function renderPlayers(players: Player[]): void {
  if (!dash) return
  if (players.length === 0) {
    dash.playersBody.replaceChildren(h('p', { class: 'muted' }, 'No players known yet.'))
    return
  }
  const rows = players
    .slice()
    .sort((a, b) => Number(b.online) - Number(a.online) || a.name.localeCompare(b.name))
    .map((p) =>
      h(
        'tr',
        {},
        h('td', {}, h('span', { class: `dot ${p.online ? 'dot-on' : 'dot-off'}` }), p.name),
        h('td', {}, roleTags(p)),
      ),
    )
  dash.playersBody.replaceChildren(
    h(
      'table',
      { class: 'players-table' },
      h('thead', {}, h('tr', {}, h('th', {}, 'Player'), h('th', {}, 'Roles'))),
      h('tbody', {}, ...rows),
    ),
  )
}

function roleTags(p: Player): HTMLElement {
  const wrap = h('span', { class: 'role-tags' })
  if (p.is_op) wrap.append(h('span', { class: 'tag tag-op' }, 'op'))
  if (p.is_whitelisted) wrap.append(h('span', { class: 'tag tag-wl' }, 'whitelist'))
  if (p.is_banned) wrap.append(h('span', { class: 'tag tag-ban' }, 'banned'))
  if (wrap.childNodes.length === 0) wrap.append(h('span', { class: 'muted' }, '—'))
  return wrap
}

// ------------------------------------------------------------------- misc

function startPolling(): void {
  stopPolling()
  pollTimer = window.setInterval(() => void refreshStatus(), 5000)
}

function stopPolling(): void {
  if (pollTimer !== undefined) {
    clearInterval(pollTimer)
    pollTimer = undefined
  }
}

function disconnect(): void {
  stopPolling()
  teardownConsole()
  clearConfig()
  dash = undefined
  renderConnect()
}

// --------------------------------------------------------------------- boot

if (getApiKey()) renderDashboard()
else renderConnect()
