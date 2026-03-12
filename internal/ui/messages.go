package ui

import "emby-tui-scrobbler/internal/emby"

type itemsMsg struct {
 items []emby.Item
}

type sessionMsg struct {
 progress float64
}