package domain

type UserInfo struct {
	Username string `json:"username,omitempty"`
	Token    string `json:"token,omitempty"`
}

type Chat struct {
	ChatID int    `json:"chatID"`
	Sender string `json:"sender"`
}
