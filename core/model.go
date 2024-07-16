package core

import (
	"loro-tui/domain"
	"loro-tui/http_client"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/charmbracelet/bubbles/textinput"
)

type View int

const (
	Login View = iota
	Chat
)

type ErrMsg struct{ error }

func (e ErrMsg) Error() string { return e.error.Error() }

type Model struct {
	UserInfo       *domain.UserInfo
	Login          LoginModel
	ChatModel      ChatModel
	Width          int
	Height         int
	HttpClient     *http_client.Client
	ServerEndpoint string
	ErrorApp       string
	Navigator      View
}

type LoginModel struct {
	focusIndex int
	inputs     []textinput.Model
}

func (m LoginModel) updateInputs(msg tea.Msg) tea.Cmd {
	cmds := make([]tea.Cmd, len(m.inputs))

	// Only text inputs with Focus() set will respond, so it's safe to simply
	// update all of them here without any further logic.
	for i := range m.inputs {
		m.inputs[i], cmds[i] = m.inputs[i].Update(msg)
	}

	return tea.Batch(cmds...)
}

type ChatModel struct {
}

func NewModel(width, height int, httpClient *http_client.Client) Model {
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
		Login:      m,
		Width:      width,
		Height:     height,
		UserInfo:   nil,
		HttpClient: httpClient,
		Navigator:  Login,
	}
}

func (m Model) Init() tea.Cmd {
	return textinput.Blink
}
