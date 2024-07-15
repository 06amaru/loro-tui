package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"golang.org/x/term"

	"loro-tui/http_client"
)

var (
	focusedStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))
	buttonStyle  = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FFF7DB")).
			Background(lipgloss.Color("#888B7E"))

	focusedButtonStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("#FFF7DB")).
				Background(lipgloss.Color("205"))

	dialogBoxStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("#874BFD")).
			Padding(1, 3).
			BorderTop(true).
			BorderLeft(true).
			BorderRight(true).
			BorderBottom(true)
	noStyle = lipgloss.NewStyle()
)

type LoginModel struct {
	focusIndex int
	inputs     []textinput.Model
}

type ChatModel struct {
}

type UserInfo struct {
}

type Model struct {
	UserInfo       *UserInfo
	Login          LoginModel
	ChatModel      ChatModel
	Width          int
	Height         int
	HttpClient     http_client.Client
	ServerEndpoint string
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
	}
}

func (m Model) Init() tea.Cmd {
	return textinput.Blink
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case ErrorMessage:
		// SET ERROR MESSAGE
		return nil, nil
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "esc":
			return m, tea.Quit
		case "tab":
			m.Login.focusIndex++
			m.Login.focusIndex = m.Login.focusIndex % 3

			cmds := make([]tea.Cmd, len(m.Login.inputs))
			for i := range m.Login.inputs {
				if i == m.Login.focusIndex {
					cmds[i] = m.Login.inputs[i].Focus()
					m.Login.inputs[i].PromptStyle = focusedStyle
					m.Login.inputs[i].TextStyle = focusedStyle
					continue
				}
				m.Login.inputs[i].Blur()
				m.Login.inputs[i].PromptStyle = noStyle
				m.Login.inputs[i].TextStyle = noStyle
			}

			return m, tea.Batch(cmds...)
		case "enter":
			if m.Login.focusIndex == 3 {
				username := m.Login.inputs[0].Value()
				password := m.Login.inputs[1].Value()
				body := http_client.RequestLogin{
					Username: username,
					Password: password,
				}

				bodyBytes, _ := json.Marshal(body)

				_, err := m.HttpClient.Post(m.ServerEndpoint+"/login", bodyBytes)
				if err != nil {
					return m, sendError(err)
				}
				// STORE USER
				// CHANGE FROM LOGIN TO CHAT
			}
		}
	}

	cmd := m.Login.updateInputs(msg)

	return m, cmd
}

func sendError(err error) tea.Cmd {
	errorMessage := ErrorMessage{err.Error()}
	return func() tea.Msg {
		return errorMessage
	}
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

func (m Model) View() string {
	var b strings.Builder
	button := &buttonStyle

	if m.Login.focusIndex == 2 {
		button = &focusedButtonStyle
	}

	ui := lipgloss.JoinVertical(lipgloss.Center,
		m.Login.inputs[0].View(),
		m.Login.inputs[1].View(),
		button.Render(" Enter "),
	)

	dialog := lipgloss.Place(m.Width, m.Height, lipgloss.Center, lipgloss.Center,
		dialogBoxStyle.Render(ui),
		//lipgloss.WithWhitespaceBackground(lipgloss.Color("#F25D94")),
	)

	b.WriteString(dialog)

	return b.String()
}

func main() {
	serverEndpoint := flag.String("server", "", "Chat server endpoint (required)")
	flag.Parse()

	if *serverEndpoint == "" {
		fmt.Println("Error: The -server flag is required")
		flag.Usage()
		os.Exit(1)
	}

	width, height, _ := term.GetSize(int(os.Stdout.Fd()))

	if _, err := tea.NewProgram(NewModel(width, height, *serverEndpoint)).Run(); err != nil {
		fmt.Printf("could not start program: %s\n", err)
		os.Exit(1)
	}
}
