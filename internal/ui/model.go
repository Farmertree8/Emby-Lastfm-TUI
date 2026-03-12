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

 list list.Model
 items []emby.Item

 stack  []string
 parent string

 nowPlaying  string
 pendingPlay *emby.Item

 progress      progress.Model
 progressRatio float64
}

func NewModel(c *emby.Client, p *player.Player) Model {

 l := list.New([]list.Item{}, list.NewDefaultDelegate(), 40, 20)

 prog := progress.New(progress.WithDefaultGradient())

 return Model{
  client: c,
  player: p,
  list: l,
  progress: prog,
 }
}

func (m Model) Init() tea.Cmd {
 return loadLibraries(m.client)
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {

 switch msg := msg.(type) {

 case itemsMsg:

  m.items = msg.items

  var li []list.Item

  for _, it := range msg.items {
   li = append(li, listItem(it))
  }

  m.list.SetItems(li)

 case tea.KeyMsg:

  // If we are waiting for a play mode choice, handle that first
  if m.pendingPlay != nil {
   switch msg.String() {
   case "n":
    item := *m.pendingPlay
    m.pendingPlay = nil
    m.nowPlaying = item.Name
    m.player.Play(item)
    return m, pollSessions(m.client)

   case "s":
    shuffled := player.Shuffle(m.items)
    m.pendingPlay = nil
    if len(shuffled) > 0 {
     m.nowPlaying = shuffled[0].Name
     go m.player.PlayMany(shuffled)
    }
    return m, pollSessions(m.client)

   case "esc":
    m.pendingPlay = nil
    return m, nil
   }
   return m, nil
  }

  switch msg.String() {

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

  case "q":

   return m, tea.Quit
  }

 case sessionMsg:

  m.progressRatio = msg.progress

  return m, pollSessions(m.client)
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

 return fmt.Sprintf(
  "%s\n\nNow Playing: %s\n%s",
  m.list.View(),
  m.nowPlaying,
  m.progress.ViewAs(m.progressRatio),
 )
}

func pollSessions(client *emby.Client) tea.Cmd {

 return tea.Tick(time.Second*10, func(t time.Time) tea.Msg {

  sessions, _ := client.GetSessions()

  for _, s := range sessions {

   if s.NowPlayingItem.Id == "" {
    continue
   }

   pos := s.PlayState.PositionTicks
   dur := s.NowPlayingItem.RunTimeTicks

   if dur == 0 {
    continue
   }

   ratio := float64(pos) / float64(dur)

   return sessionMsg{ratio}
  }

  return sessionMsg{0}
 })
}