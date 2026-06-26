package types

type Player struct {
	UUID          string `json:"uuid"`
	Name          string `json:"name"`
	Online        bool   `json:"online"`
	IsOp          bool   `json:"is_op"`
	IsBanned      bool   `json:"is_banned"`
	IsWhitelisted bool   `json:"is_whitelisted"`
}

type UserCacheEntry struct {
	UUID      string `json:"uuid"`
	Name      string `json:"name"`
	ExpiresOn string `json:"expires_on"`
}

type OpEntry struct {
	UUID                string `json:"uuid"`
	Name                string `json:"name"`
	Level               int    `json:"level"`
	BypassesPlayerLimit bool   `json:"bypasses_player_limit"`
}

type BannedPlayerEntry struct {
	UUID    string `json:"uuid"`
	Name    string `json:"name"`
	Created string `json:"created"`
	Source  string `json:"source"`
	Expires string `json:"expires"`
	Reason  string `json:"reason"`
}

type WhitelistEntry struct {
	UUID string `json:"uuid"`
	Name string `json:"name"`
}
