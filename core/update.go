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
					return m, func() tea.Msg { return LoadingEvent{} }
				}
			}
		}

		cmd := m.Login.updateInputs(msg)

		return m, cmd
	}
	if m.Navigator == Chat {
		switch msg := msg.(type) {
		case LoadingEvent:
			chats, err := m.HttpClient.GetChats(m.UserInfo.Token)
			if err != nil {
				return m, func() tea.Msg { return err }
			}

			items := []list.Item{}
			for _, chat := range chats {
				items = append(items, widgets.NewItem(chat.Sender, chat.ChatID))
			}

			m.ChatModel.List.SetItems(items)
			if len(items) != 0 {
				defaultChat := m.ChatModel.List.SelectedItem().(widgets.Item)
				url := fmt.Sprintf("ws://localhost:8081/socket/join?id=%d", defaultChat.ChatID)
				newSocket, err := web_socket.NewWSocketClient(url, m.UserInfo.Token)
				if err != nil {
					return m, func() tea.Msg { return err }
				}
				if m.Socket != nil {
					m.Socket.Close()
				}
				m.Socket = newSocket

				go func(model *Model) {
					for {
						bytes, err := m.Socket.Listen()
						if err != nil {
							m.Program.Send(err)
							break
						}
						m.Program.Send(SocketMsg{bytes})
					}
				}(&m)

				return m, nil
			}
		case error:
			m.ErrorApp = msg.Error()
			return m, nil
		case tea.KeyMsg:
			switch msg.String() {
			case "tab":
				m.ChatModel.focusIndex++
				m.ChatModel.focusIndex = m.ChatModel.focusIndex % 3
				if m.ChatModel.focusIndex == 0 {
					m.ChatModel.Input.TextStyle = focusedStyle
					m.ChatModel.Input.PromptStyle = focusedStyle
					m.ChatModel.Input.Cursor.Style = focusedStyle
				} else {
					m.ChatModel.Input.TextStyle = noStyle
					m.ChatModel.Input.PromptStyle = noStyle
					m.ChatModel.Input.Cursor.Style = noStyle
				}

				if m.ChatModel.focusIndex == 1 {
					delegateList := widgets.NewDelegateList(true)
					m.ChatModel.List.SetDelegate(delegateList)
				} else {
					delegateList := widgets.NewDelegateList(false)
					m.ChatModel.List.SetDelegate(delegateList)
				}
				return m, nil
			case "ctrl+c", "esc":
				return m, tea.Quit
			case "ctrl+n":
				m.Navigator = NewChat
				return m, nil
			case "enter":
				if m.ChatModel.focusIndex == 0 && m.Socket != nil {
					message := m.ChatModel.Input.Value()
					err := m.Socket.Send(message)
					if err != nil {
						return m, func() tea.Msg { return err }
					}
					m.ChatModel.Input.SetValue("")
					return m, nil
				}

				if m.ChatModel.focusIndex == 1 && len(m.ChatModel.List.Items()) != 0 {
					defaultChat := m.ChatModel.List.SelectedItem().(widgets.Item)
					url := fmt.Sprintf("ws://localhost:8081/socket/join?id=%d", defaultChat.ChatID)
					newSocket, err := web_socket.NewWSocketClient(url, m.UserInfo.Token)
					if err != nil {
						return m, func() tea.Msg { return err }
					}
					if m.Socket != nil {
						m.Socket.Close()
					}
					m.Socket = newSocket

					go func(model *Model) {
						for {
							bytes, err := m.Socket.Listen()
							if err != nil {
								m.Program.Send(err)
								break
							}
							m.Program.Send(SocketMsg{bytes})
						}

					}(&m)

					return m, nil
				}
			}

			var cmd tea.Cmd
			if m.ChatModel.focusIndex == 0 {
				m.ChatModel.Input, cmd = m.ChatModel.Input.Update(msg)
				return m, cmd
			}
			if m.ChatModel.focusIndex == 1 {
				m.ChatModel.List, cmd = m.ChatModel.List.Update(msg)
				return m, cmd
			}

			return m, nil
		case SocketMsg:
			message := string(msg.Data)
			m.ChatModel.Messages += fmt.Sprintf("%s\n", message)
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
					_, err := web_socket.NewWSocketClient(url, m.UserInfo.Token)
					if err != nil {
						return m, func() tea.Msg { return err }
					}

					m.Navigator = Chat
					return m, func() tea.Msg { return LoadingEvent{} }
				}
			}
			var cmd tea.Cmd
			m.NewChat.Input, cmd = m.NewChat.Input.Update(msg)
			return m, cmd
		}
	}
	return m, nil
}
