package http_client

type RequestLogin struct {
	Username string `json:"username"`
	Password string `json:"password"`
}
