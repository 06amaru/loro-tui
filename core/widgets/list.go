package widgets

import "github.com/charmbracelet/bubbles/list"

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
	delegateList := list.NewDefaultDelegate()
	delegateList.ShowDescription = false

	chatList := list.New(items, delegateList, width, height)

	return chatList
}
