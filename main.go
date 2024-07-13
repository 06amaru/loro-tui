package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"golang.org/x/term"
)

var (
	focusedStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))
	buttonStyle  = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FFF7DB")).
			Background(lipgloss.Color("#888B7E"))

	focusedButtonStyle = buttonStyle.
				Foreground(lipgloss.Color("#FFF7DB")).
				Background(lipgloss.Color("#F25D94"))

	dialogBoxStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("#874BFD")).
			Padding(1, 3).
			BorderTop(true).
			BorderLeft(true).
			BorderRight(true).
			BorderBottom(true)
)

type LoginModel struct {
	focusIndex int
	inputs     []textinput.Model
}

type Model struct {
	Login  LoginModel
	Width  int
	Height int
}

func NewModel(width, height int) Model {
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
		Login:  m,
		Width:  width,
		Height: height,
	}
}

func (m Model) Init() tea.Cmd {
	return textinput.Blink
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "esc":
			return m, tea.Quit
		}
	}

	return m, nil
}

func (m Model) View() string {
	var b strings.Builder
	button := &buttonStyle

	ui := lipgloss.JoinVertical(lipgloss.Center,
		m.Login.inputs[0].View(),
		m.Login.inputs[1].View(),
		button.Render("ENTER"),
	)

	dialog := lipgloss.Place(m.Width, m.Height, lipgloss.Center, lipgloss.Center,
		dialogBoxStyle.Render(ui),
		//lipgloss.WithWhitespaceBackground(lipgloss.Color("#F25D94")),
	)

	b.WriteString(dialog)

	return b.String()
}

func main() {
	width, height, _ := term.GetSize(int(os.Stdout.Fd()))

	if _, err := tea.NewProgram(NewModel(width, height)).Run(); err != nil {
		fmt.Printf("could not start program: %s\n", err)
		os.Exit(1)
	}
}
