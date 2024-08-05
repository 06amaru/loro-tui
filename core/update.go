package core

import (
	"loro-tui/core/widgets"
	"loro-tui/http_client"

	"github.com/charmbracelet/bubbles/list"
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
					//TODO: LOAD CHATS
					items := []list.Item{
						widgets.NewItem("godwana", 1),
						widgets.NewItem("morodo", 2),
						widgets.NewItem("eminem", 3),
						widgets.NewItem("50cent", 4),
						widgets.NewItem("snoopdog", 5),
					}
					m.Chat.List.SetItems(items)
					//m.Chat.List.SetSize(m.Width, m.Height)
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
			case "tab":
				m.Chat.focusIndex++
				m.Chat.focusIndex = m.Chat.focusIndex % 3
				return m, nil
			case "ctrl+c", "esc":
				return m, tea.Quit
			case "ctrl+n":
				m.Navigator = NewChat
				return m, nil
			}

			var cmd tea.Cmd
			if m.Chat.focusIndex == 0 {
				m.Chat.Input, cmd = m.Chat.Input.Update(msg)
			}
			if m.Chat.focusIndex == 1 {
				m.Chat.List, cmd = m.Chat.List.Update(msg)
			}

			return m, cmd
		}
	}
	return m, nil
}
