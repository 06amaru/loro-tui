package web_socket

import (
	"fmt"
	"net/http"

	"github.com/gorilla/websocket"
)

type WSocketClient struct {
	conn *websocket.Conn
}

func NewWSocketClient(url, token string) (*WSocketClient, error) {

	header := http.Header{}
	header.Set("Cookie", fmt.Sprintf("token=%s", token))

	conn, _, err := websocket.DefaultDialer.Dial(url, header)
	if err != nil {
		return nil, fmt.Errorf("websocket new client [%v]", err)
	}

	return &WSocketClient{conn: conn}, nil
}

func (ws WSocketClient) Send(message string) error {
	err := ws.conn.WriteMessage(websocket.TextMessage, []byte(message))
	if err != nil {
		return fmt.Errorf("websocket send [%v]", err)
	}
	return nil
}

func (ws WSocketClient) Listen() ([]byte, error) {
	_, message, err := ws.conn.ReadMessage()
	if err != nil {
		return nil, fmt.Errorf("websocket listen [%v]", err)
	}

	return message, nil
}

func (ws WSocketClient) Close() {
	ws.conn.Close()
}
