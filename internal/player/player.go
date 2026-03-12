package player

import (
	"log"
	"math/rand"
	"os"
	"os/exec"
	"time"

	"emby-tui-scrobbler/internal/emby"
	"emby-tui-scrobbler/internal/lastfm"
)

type Player struct {
	emby   *emby.Client
	lastfm *lastfm.Client
}

func New(e *emby.Client, l *lastfm.Client) *Player {
	return &Player{
		emby:   e,
		lastfm: l,
	}
}

func (p *Player) Play(item emby.Item) {

 url := p.emby.StreamURL(item.Id)

 cmd := exec.Command("mpv", "--no-video", url)

 cmd.Stdout = os.Stdout
 cmd.Stderr = os.Stderr

 err := cmd.Start()
 if err != nil {
  log.Println("mpv failed:", err)
  return
 }

 artist := ""
 if len(item.Artists) > 0 {
  artist = item.Artists[0]
 }

 start := time.Now().Unix()
 duration := time.Duration(item.RunTimeTicks/10000000) * time.Second

 p.lastfm.UpdateNowPlaying(item.Name, artist)

 go func() {
  time.Sleep(duration / 2)
  p.lastfm.Scrobble(item.Name, artist, start)
 }()

 go cmd.Wait()
}

func (p *Player) playAndWait(item emby.Item) {

 url := p.emby.StreamURL(item.Id)

 cmd := exec.Command("mpv", "--no-video", url)

 cmd.Stdout = os.Stdout
 cmd.Stderr = os.Stderr

 artist := ""
 if len(item.Artists) > 0 {
  artist = item.Artists[0]
 }

 start := time.Now().Unix()
 duration := time.Duration(item.RunTimeTicks/10000000) * time.Second

 p.lastfm.UpdateNowPlaying(item.Name, artist)

 go func() {
  time.Sleep(duration / 2)
  p.lastfm.Scrobble(item.Name, artist, start)
 }()

 // Blocks until mpv exits — no Sleep needed
 if err := cmd.Run(); err != nil {
  log.Println("mpv error:", err)
 }
}

func (p *Player) PlayMany(items []emby.Item) {

 for _, it := range items {
  p.playAndWait(it)
 }
}

func Shuffle(items []emby.Item) []emby.Item {

	rand.Seed(time.Now().UnixNano())

	rand.Shuffle(len(items), func(i, j int) {
		items[i], items[j] = items[j], items[i]
	})

	return items
}