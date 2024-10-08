package core

import "github.com/charmbracelet/lipgloss"

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
	noStyle        = lipgloss.NewStyle()
	statusBarStyle = lipgloss.NewStyle().
			Foreground(lipgloss.AdaptiveColor{Light: "#343433", Dark: "#C1C6B2"}).
			Background(lipgloss.AdaptiveColor{Light: "#D9DCCF", Dark: "#353533"})
	statusStyle = lipgloss.NewStyle().
			Inherit(statusBarStyle).
			Foreground(lipgloss.Color("#FFFDF5")).
			Background(lipgloss.Color("#FF5F87")).
			Padding(0, 1).
			MarginRight(1)
	statusText = lipgloss.NewStyle().Inherit(statusBarStyle)
)
