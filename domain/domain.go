package domain

type UserInfo struct {
	Username string `json:"username,omitempty"`
	Token    string `json:"token,omitempty"`
}
