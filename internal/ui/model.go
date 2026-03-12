package ui

import (
	"fmt"
	"time"

	"emby-tui-scrobbler/internal/emby"
	"emby-tui-scrobbler/internal/player"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/progress"
	tea "github.com/charmbracelet/bubbletea"
)

type Model struct {
	client *emby.Client
	player *player.Player

	list  list.Model
	items []emby.Item

	stack  []string
	parent string

	nowPlaying    string
	pendingPlay   *emby.Item

	progress      progress.Model
	progressRatio float64

	// Local playback clock — set when a track starts, zeroed on stop.
	trackStart    time.Time
	trackDuration time.Duration // 0 means nothing is playing

	program *tea.Program
}

func NewModel(c *emby.Client, p *player.Player) Model {

	l := list.New([]list.Item{}, list.NewDefaultDelegate(), 40, 20)
	prog := progress.New(progress.WithDefaultGradient())

	return Model{
		client:   c,
		player:   p,
		list:     l,
		progress: prog,
	}
}

func (m *Model) SetProgram(p *tea.Program) {
	m.program = p
	m.player.OnTrackChange = func(name string, durationTicks int64) {
		p.Send(trackChangeMsg{name, durationTicks})
	}
}

func (m Model) Init() tea.Cmd {
	return loadLibraries(m.client)
}

// tick fires once per second while a track is playing.
func tick() tea.Cmd {
	return tea.Tick(time.Second, func(time.Time) tea.Msg {
		return tickMsg{}
	})
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {

	switch msg := msg.(type) {

	case tea.WindowSizeMsg:
		m.progress.Width = msg.Width
		m.list.SetSize(msg.Width, msg.Height-5)
		return m, nil

	case trackChangeMsg:
		m.nowPlaying = msg.name
		m.progressRatio = 0
		if msg.name == "" || msg.durationTicks == 0 {
			// Playback stopped or unknown duration — kill the ticker.
			m.trackDuration = 0
			return m, nil
		}
		m.trackStart = time.Now()
		m.trackDuration = time.Duration(msg.durationTicks/10000000) * time.Second
		return m, tick()

	case tickMsg:
		if m.trackDuration == 0 {
			// Nothing playing; let the ticker die.
			return m, nil
		}
		elapsed := time.Since(m.trackStart)
		m.progressRatio = elapsed.Seconds() / m.trackDuration.Seconds()
		if m.progressRatio > 1 {
			m.progressRatio = 1
		}
		return m, tick()

	case itemsMsg:
		m.items = msg.items
		var li []list.Item
		for _, it := range msg.items {
			li = append(li, listItem(it))
		}
		m.list.SetItems(li)
		return m, nil

	case tea.KeyMsg:

		// When waiting for play mode choice, consume ALL keys here.
		if m.pendingPlay != nil {
			switch msg.String() {
			case "n":
				selected := *m.pendingPlay
				m.pendingPlay = nil

				queue := []emby.Item{}
				found := false
				for _, it := range m.items {
					if it.Id == selected.Id {
						found = true
					}
					if found {
						queue = append(queue, it)
					}
				}
				if len(queue) == 0 {
					queue = []emby.Item{selected}
				}

				m.nowPlaying = selected.Name
				go m.player.PlayMany(queue)
				return m, nil

			case "s":
				selected := *m.pendingPlay
				rest := make([]emby.Item, 0, len(m.items)-1)
				for _, it := range m.items {
					if it.Id != selected.Id {
						rest = append(rest, it)
					}
				}
				rest = player.Shuffle(rest)
				queue := append([]emby.Item{selected}, rest...)
				m.pendingPlay = nil
				m.nowPlaying = selected.Name
				go m.player.PlayMany(queue)
				return m, nil

			case "esc":
				m.pendingPlay = nil
				return m, nil
			}
			return m, nil
		}

		switch msg.String() {

		case "q":
			return m, tea.Quit

		case "s":
			m.player.Skip()
			return m, nil

		case "delete":
			m.player.Stop()
			m.nowPlaying = ""
			m.progressRatio = 0
			m.trackDuration = 0
			return m, nil

		case "backspace":
			if len(m.stack) == 0 {
				return m, loadLibraries(m.client)
			}
			last := m.stack[len(m.stack)-1]
			m.stack = m.stack[:len(m.stack)-1]
			m.parent = last
			if last == "" {
				return m, loadLibraries(m.client)
			}
			return m, loadItems(m.client, last)

		case "enter":
			idx := m.list.Index()
			if idx >= len(m.items) {
				break
			}
			item := m.items[idx]
			if item.Type == "Audio" {
				m.pendingPlay = &item
				return m, nil
			}
			m.stack = append(m.stack, m.parent)
			m.parent = item.Id
			return m, loadItems(m.client, item.Id)
		}
	}

	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

func (m Model) View() string {

	if m.pendingPlay != nil {
		return fmt.Sprintf(
			"%s\n\nPlay \"%s\": [n] Normal  [s] Shuffle all  [esc] Cancel",
			m.list.View(),
			m.pendingPlay.Name,
		)
	}

	nowPlayingLine := "Now Playing: " + m.nowPlaying
	if m.nowPlaying == "" {
		nowPlayingLine = "Now Playing: —"
	}

	elapsed := ""
	if m.trackDuration > 0 {
		e := time.Since(m.trackStart).Truncate(time.Second)
		d := m.trackDuration.Truncate(time.Second)
		elapsed = fmt.Sprintf(" %d:%02d / %d:%02d",
			int(e.Minutes()), int(e.Seconds())%60,
			int(d.Minutes()), int(d.Seconds())%60,
		)
	}

	return fmt.Sprintf(
		"%s\n\n%s%s\n%s\n[s] Skip  [del] Stop  [q] Quit",
		m.list.View(),
		nowPlayingLine,
		elapsed,
		m.progress.ViewAs(m.progressRatio),
	)
}
