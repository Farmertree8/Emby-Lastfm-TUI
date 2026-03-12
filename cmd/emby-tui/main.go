package main

import (
 "log"

 "emby-tui-scrobbler/internal/config"
 "emby-tui-scrobbler/internal/emby"
 "emby-tui-scrobbler/internal/lastfm"
 "emby-tui-scrobbler/internal/player"
 "emby-tui-scrobbler/internal/ui"

 tea "github.com/charmbracelet/bubbletea"
)

func main() {

 cfg, err := config.Load("config.json")
 if err != nil {
  log.Fatal(err)
 }

 embyClient := emby.New(
  cfg.EmbyURL,
  cfg.EmbyAPIKey,
  cfg.UserID,
 )

 lastfmClient := lastfm.New(
  cfg.LastFMApiKey,
  cfg.LastFMSecret,
  cfg.LastFMSessionKey,
 )

 player := player.New(embyClient, lastfmClient)

 model := ui.NewModel(embyClient, player)

 program := tea.NewProgram(model, tea.WithAltScreen())

 if err := program.Start(); err != nil {
  log.Fatal(err)
 }
}