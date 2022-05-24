package types

// ServerInfo necessary fields returned from platform
type ServerInfo struct {
	CookieName   string `json:"cookieName"`
	SecureCookie bool   `json:"secureCookie"`
}

type Result struct {
	ID       string `json:"_id"`
	Name     string `json:"name,omitempty"`
	Username string `json:"username,omitempty"`
}

type ResultFilter struct {
	Result      []Result `json:"result"`
	ResultCount int      `json:"resultCount,omitempty"`
}
