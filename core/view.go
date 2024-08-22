package core

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
)

func (m Model) View() string {
	var b strings.Builder

	if m.Navigator == Login {
		button := &buttonStyle

		if m.Login.focusIndex == 2 {
			button = &focusedButtonStyle
		}

		ui := lipgloss.JoinVertical(lipgloss.Center,
			m.Login.inputs[0].View(), // username
			m.Login.inputs[1].View(), // password
			button.Render(" Enter "),
		)

		dialog := lipgloss.Place(m.Width, m.Height, lipgloss.Center, lipgloss.Center,
			dialogBoxStyle.Render(ui),
		)

		w := lipgloss.Width
		statusKey := statusStyle.Render("STATUS")
		statusVal := statusText.
			Width(m.Width - w(statusKey)).Render(m.ErrorApp)

		bar := lipgloss.JoinHorizontal(lipgloss.Top, statusKey, statusVal)

		dialog = lipgloss.JoinVertical(lipgloss.Bottom, dialog, bar)

		b.WriteString(dialog)
	}

	if m.Navigator == NewChat {
		button := &buttonStyle
		if m.NewChat.focusIndex == 1 {
			button = &focusedButtonStyle
		}

		ui := lipgloss.JoinVertical(lipgloss.Center,
			m.NewChat.Input.View(),
			button.Render(" Chat "),
		)

		dialog := lipgloss.Place(m.Width, m.Height, lipgloss.Center, lipgloss.Center,
			dialogBoxStyle.Render(ui),
		)

		w := lipgloss.Width
		statusKey := statusStyle.Render("STATUS")
		statusVal := statusText.
			Width(m.Width - w(statusKey)).Render(m.ErrorApp)

		bar := lipgloss.JoinHorizontal(lipgloss.Top, statusKey, statusVal)

		dialog = lipgloss.JoinVertical(lipgloss.Bottom, dialog, bar)

		b.WriteString(dialog)
	}

	if m.Navigator == Chat {
		if m.ChatModel.Messages != "" {
			m.ChatModel.Chat.SetContent(m.ChatModel.Messages)
		}

		var ui string
		m.ChatModel.Chat.Height = m.Height - 1

		ui = lipgloss.JoinHorizontal(lipgloss.Center,
			m.ChatModel.Chat.View(),
			m.ChatModel.List.View(),
		)
		ui = lipgloss.JoinVertical(lipgloss.Top,
			ui,
			m.ChatModel.Input.View(), // body message
		)
		b.WriteString(ui)
	}

	return b.String()
}
