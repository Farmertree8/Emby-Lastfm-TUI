package ui

import (
 "emby-tui-scrobbler/internal/emby"

 tea "github.com/charmbracelet/bubbletea"
)

func loadLibraries(client *emby.Client) tea.Cmd {

 return func() tea.Msg {

  items, _ := client.GetLibraries()

  return itemsMsg{items}
 }
}

func loadItems(client *emby.Client, parent string) tea.Cmd {

 return func() tea.Msg {

  items, _ := client.GetItems(parent)

  return itemsMsg{items}
 }
}