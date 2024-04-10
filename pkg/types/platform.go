package types

type JWTPayload struct {
	ISS string `json:"iss"`
	SUB string `json:"sub"`
	AUD string `json:"aud"`
	JTI string `json:"jti"`
	EXP int64  `json:"exp"`
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

// ServerInfo necessary fields returned from platform
type ServerInfo struct {
	CookieName   string `json:"cookieName"`
	SecureCookie bool   `json:"secureCookie"`
}