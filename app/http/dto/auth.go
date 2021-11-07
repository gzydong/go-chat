package dto

type TokenAttribute struct {
	Type      string `json:"type"`
	Token     string `json:"access_token"`
	ExpiresIn int64  `json:"expires_in"`
}
