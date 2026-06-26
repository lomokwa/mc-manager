# MC Manager — Web UI

A small **Vite + TypeScript** (no framework) frontend for the mc-manager API:
server status, start/stop, a live console over WebSocket, and the player list.

## Develop

```bash
cd web
npm install
npm run dev
```

The dev server runs on <http://localhost:5173> and proxies `/api` (REST plus
the `/api/console` WebSocket) to `http://localhost:8080`. Point it elsewhere
with:

```bash
VITE_API_PROXY=http://my-host:8080 npm run dev
```

On first load, enter your `API_KEY` (and optionally the API URL — leave it
blank to use the same origin, which is what the dev proxy expects). Settings
are stored in `localStorage`; use **Change key** to reset them.

## Build

```bash
npm run build      # type-checks, then bundles to web/dist
npm run preview    # serve the production build locally
```

The static files in `dist/` can be served by any static host — or by the Go
server itself.
