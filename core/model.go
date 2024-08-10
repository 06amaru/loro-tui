package core

import (
	"loro-tui/core/widgets"
	"loro-tui/domain"
	"loro-tui/http_client"
	"loro-tui/web_socket"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type View int

const (
	Login View = iota
	Chat
	NewChat
)

type ErrMsg struct{ error }

func (e ErrMsg) Error() string { return e.error.Error() }

type SocketMsg struct {
	Data []byte
}

type Model struct {
	UserInfo   *domain.UserInfo
	Login      LoginModel
	Chat       ChatModel
	NewChat    NewChatModel
	Width      int
	Height     int
	HttpClient *http_client.Client
	ErrorApp   string
	Navigator  View
	Socket     *web_socket.WSocketClient
}

type NewChatModel struct {
	focusIndex int
	Input      textinput.Model
}

type LoginModel struct {
	focusIndex int
	inputs     []textinput.Model
}

type ChatModel struct {
	List       list.Model
	Message    viewport.Model
	focusIndex int
	Input      textinput.Model
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

func NewModel(width, height int, httpClient *http_client.Client) Model {
	loginM := LoginModel{
		inputs: make([]textinput.Model, 2),
	}

	var t textinput.Model
	for i := range loginM.inputs {
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

		loginM.inputs[i] = t
	}

	chatList := widgets.NewList(nil, width, height)
	chatM := ChatModel{
		List:    chatList,
		Message: viewport.New(width-(width/3), height),
	}

	chatM.Message.Style = lipgloss.NewStyle().
		BorderStyle(lipgloss.NormalBorder()).
		MarginRight(4).
		MarginTop(2)

	t = textinput.New()
	t.CharLimit = 32
	t.Width = 33
	t.Placeholder = "Write your message"
	t.TextStyle = focusedStyle
	t.PromptStyle = focusedStyle
	t.Cursor.Style = focusedStyle
	t.Focus()
	chatM.Input = t

	newChatM := NewChatModel{}
	t = textinput.New()
	t.CharLimit = 32
	t.Width = 33
	t.Placeholder = "Here goes the username to chat"
	t.TextStyle = focusedStyle
	t.PromptStyle = focusedStyle
	t.Cursor.Style = focusedStyle
	t.Focus()
	newChatM.Input = t

	return Model{
		Login:      loginM,
		Chat:       chatM,
		NewChat:    newChatM,
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
