package ui

import "emby-tui-scrobbler/internal/emby"

type itemsMsg struct {
	items []emby.Item
}

type tickMsg struct{}

type trackChangeMsg struct {
	name         string
	durationTicks int64 // RunTimeTicks; 0 when playback stops
}
