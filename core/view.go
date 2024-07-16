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
			m.Login.inputs[0].View(),
			m.Login.inputs[1].View(),
			button.Render(" Enter "),
		)

		dialog := lipgloss.Place(m.Width, m.Height, lipgloss.Center, lipgloss.Center,
			dialogBoxStyle.Render(ui),
			//lipgloss.WithWhitespaceBackground(lipgloss.Color("#F25D94")),
		)

		w := lipgloss.Width
		statusKey := statusStyle.Render("STATUS")
		statusVal := statusText.
			Width(m.Width - w(statusKey)).Render(m.ErrorApp)

		bar := lipgloss.JoinHorizontal(lipgloss.Top, statusKey, statusVal)

		//b.WriteString(statusBarStyle.Width(m.Width).Render(bar))

		dialog = lipgloss.JoinVertical(lipgloss.Bottom, dialog, bar)

		b.WriteString(dialog)
	}

	if m.Navigator == Chat {
		b.WriteString("Hi " + m.UserInfo.Username)
	}

	return b.String()
}
