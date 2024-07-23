package core

import (
	"loro-tui/http_client"

	tea "github.com/charmbracelet/bubbletea"
)

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if m.Navigator == Login {
		switch msg := msg.(type) {
		case error:
			m.ErrorApp = msg.Error()
			return m, nil
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
				if m.Login.focusIndex == 2 {
					m.ErrorApp = ""
					username := m.Login.inputs[0].Value()
					password := m.Login.inputs[1].Value()
					body := http_client.RequestLogin{
						Username: username,
						Password: password,
					}

					user, err := m.HttpClient.Login(body)
					if err != nil {
						return m, func() tea.Msg { return err }
					}
					m.UserInfo = user
					m.Navigator = Chat
					return m, nil
				}
			}
		}

		cmd := m.Login.updateInputs(msg)

		return m, cmd
	}
	if m.Navigator == Chat {
		switch msg := msg.(type) {
		case error:
			m.ErrorApp = msg.Error()
			return m, nil
		case tea.KeyMsg:
			switch msg.String() {
			case "ctrl+c", "esc":
				return m, tea.Quit
			}
		}
		return m, nil
	}

	return m, nil
}
