package core

import (
	"loro-tui/http_client"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/charmbracelet/bubbles/textinput"
)

type ErrMsg struct{ error }

func (e ErrMsg) Error() string { return e.error.Error() }

type Model struct {
	UserInfo       *UserInfo
	Login          LoginModel
	ChatModel      ChatModel
	Width          int
	Height         int
	HttpClient     *http_client.Client
	ServerEndpoint string
	ErrorApp       string
}

type LoginModel struct {
	focusIndex int
	inputs     []textinput.Model
}

type ChatModel struct {
}

type UserInfo struct {
}

type ErrorMessage struct {
	Data string
}

func NewModel(width, height int, serverEndpoint string) Model {
	m := LoginModel{
		inputs: make([]textinput.Model, 2),
	}

	var t textinput.Model
	for i := range m.inputs {
		t = textinput.New()
		t.Cursor.Style = focusedStyle
		t.CharLimit = 32
		t.Width = 33

		switch i {
		case 0:
			t.Placeholder = "Username"
			t.Focus()
			t.PromptStyle = focusedStyle
			t.TextStyle = focusedStyle
		case 1:
			t.Placeholder = "Password"
			t.EchoMode = textinput.EchoPassword
			t.EchoCharacter = '*'
		}

		m.inputs[i] = t
	}

	return Model{
		Login:          m,
		Width:          width,
		Height:         height,
		UserInfo:       nil,
		ServerEndpoint: serverEndpoint,
		HttpClient:     http_client.NewClient(),
	}
}

func (m Model) Init() tea.Cmd {
	return textinput.Blink
}
