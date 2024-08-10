package widgets

import (
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/lipgloss"
)

type Item struct {
	Sender string
	ChatID int
}

func (i Item) FilterValue() string {
	return i.Sender
}

func (i Item) Title() string {
	return i.Sender
}

func (i Item) Description() string {
	return ""
}

func NewItem(sender string, chatID int) Item {
	return Item{
		Sender: sender,
		ChatID: chatID,
	}
}

func NewList(items []list.Item, width, height int) list.Model {
	delegateList := NewDelegateList(false)
	chatList := list.New(items, delegateList, width, height)
	return chatList
}

func NewDelegateList(focus bool) list.ItemDelegate {
	delegateList := list.NewDefaultDelegate()
	delegateList.ShowDescription = false
	if !focus {
		delegateList.Styles.SelectedTitle = lipgloss.NewStyle()
	}
	return delegateList
}
