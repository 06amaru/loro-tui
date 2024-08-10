package core

import (
	"fmt"
	"loro-tui/core/widgets"
	"loro-tui/http_client"
	"loro-tui/web_socket"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
)

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	defer func() {
		if m.Socket != nil {
			m.Socket.Close()
		}
	}()

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
						widgets.NewItem("godwana", 121),
						widgets.NewItem("morodo", 242),
						widgets.NewItem("eminem", 364),
						widgets.NewItem("50cent", 487),
						widgets.NewItem("snoopdog", 598),
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
				if m.Chat.focusIndex == 0 {
					m.Chat.Input.TextStyle = focusedStyle
					m.Chat.Input.PromptStyle = focusedStyle
					m.Chat.Input.Cursor.Style = focusedStyle
				} else {
					m.Chat.Input.TextStyle = noStyle
					m.Chat.Input.PromptStyle = noStyle
					m.Chat.Input.Cursor.Style = noStyle
				}

				if m.Chat.focusIndex == 1 {
					delegateList := widgets.NewDelegateList(true)
					m.Chat.List.SetDelegate(delegateList)
				} else {
					delegateList := widgets.NewDelegateList(false)
					m.Chat.List.SetDelegate(delegateList)
				}
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
				return m, cmd
			}
			if m.Chat.focusIndex == 1 {
				m.Chat.List, cmd = m.Chat.List.Update(msg)
				return m, cmd
			}

			return m, nil
		case SocketMsg:

		}
	}
	if m.Navigator == NewChat {
		switch msg := msg.(type) {
		case error:
			m.ErrorApp = msg.Error()
			return m, nil
		case tea.KeyMsg:
			switch msg.String() {
			case "ctrl+c", "esc":
				m.ErrorApp = ""
				m.Navigator = Chat
				return m, nil
			case "tab":
				m.NewChat.focusIndex++
				m.NewChat.focusIndex = m.NewChat.focusIndex % 2

				var cmd tea.Cmd
				if m.NewChat.focusIndex == 0 {
					cmd = m.NewChat.Input.Focus()
					m.NewChat.Input.PromptStyle = focusedStyle
					m.NewChat.Input.TextStyle = focusedStyle
				} else {
					m.NewChat.Input.Blur()
					m.NewChat.Input.PromptStyle = noStyle
					m.NewChat.Input.TextStyle = noStyle
				}

				return m, cmd
			case "enter":
				if m.NewChat.focusIndex == 1 {
					username := m.NewChat.Input.Value()
					url := fmt.Sprintf("ws://localhost:8081/socket/create-chat?to=%s", username)
					newSocket, err := web_socket.NewWSocketClient(url, m.UserInfo.Token)
					if err != nil {
						return m, func() tea.Msg { return err }
					}
					// close previous connection
					if m.Socket != nil {
						m.Socket.Close()
					}
					m.Socket = newSocket
					go func(model *Model) {
						defer func() {
							if m.Socket != nil {
								m.Socket.Close()
							}
						}()

						for {
							bytes, err := m.Socket.Listen()
							if err != nil {
								m.Update(func() tea.Msg { return err })
								break
							}
							m.Update(func() tea.Msg { return SocketMsg{bytes} })
						}
					}(&m)
				}
			}
			var cmd tea.Cmd
			m.NewChat.Input, cmd = m.NewChat.Input.Update(msg)
			return m, cmd
		}
	}
	return m, nil
}
