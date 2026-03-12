package player

import (
	"fmt"
	"log"
	"math/rand"
	"os"
	"os/exec"
	"sync"
	"syscall"
	"time"

	"emby-tui-scrobbler/internal/emby"
	"emby-tui-scrobbler/internal/lastfm"
)

type Player struct {
	emby   *emby.Client
	lastfm *lastfm.Client

	mu      sync.Mutex
	current *os.Process
	skip    chan struct{}
	stop    chan struct{}
	stopped bool

	// name is empty and durationTicks is 0 when playback stops.
	OnTrackChange func(name string, durationTicks int64)
}

func New(e *emby.Client, l *lastfm.Client) *Player {
	return &Player{
		emby:   e,
		lastfm: l,
	}
}

func (p *Player) Play(item emby.Item) {

	url := p.emby.StreamURL(item.Id)

	cmd := exec.Command("mpv", "--no-video", "--cover-art-auto=fuzzy", url)

	cmd.Stdout = nil
	cmd.Stderr = nil
	cmd.SysProcAttr = &syscall.SysProcAttr{CreationFlags: syscall.CREATE_NEW_PROCESS_GROUP}

	err := cmd.Start()
	if err != nil {
		log.Println("mpv failed:", err)
		return
	}

	p.mu.Lock()
	p.current = cmd.Process
	p.mu.Unlock()

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

	cmd := exec.Command("mpv", "--no-video", "--cover-art-auto=fuzzy", url)

	cmd.Stdout = nil
	cmd.Stderr = nil
	cmd.SysProcAttr = &syscall.SysProcAttr{CreationFlags: syscall.CREATE_NEW_PROCESS_GROUP}

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

	if err := cmd.Start(); err != nil {
		log.Println("mpv error:", err)
		return
	}

	p.mu.Lock()
	p.current = cmd.Process
	p.mu.Unlock()

	done := make(chan error, 1)
	go func() {
		done <- cmd.Wait()
	}()

	killMpv := func() {
		kill := exec.Command("taskkill", "/F", "/T", "/PID", fmt.Sprintf("%d", cmd.Process.Pid))
		kill.Run()
		<-done
	}

	select {
	case <-done:
		// finished naturally
	case <-p.skip:
		killMpv()
	case <-p.stop:
		killMpv()
	}
}

func (p *Player) PlayMany(items []emby.Item) {

	p.mu.Lock()
	p.skip = make(chan struct{}, 1)
	p.stop = make(chan struct{})
	p.stopped = false
	p.mu.Unlock()

	for _, it := range items {

		p.mu.Lock()
		stopped := p.stopped
		p.mu.Unlock()

		if stopped {
			if p.OnTrackChange != nil {
				p.OnTrackChange("", 0)
			}
			return
		}

		if p.OnTrackChange != nil {
			p.OnTrackChange(it.Name, it.RunTimeTicks)
		}

		p.playAndWait(it)
	}

	// Queue finished naturally
	if p.OnTrackChange != nil {
		p.OnTrackChange("", 0)
	}

	p.mu.Lock()
	p.current = nil
	p.mu.Unlock()
}

func (p *Player) Skip() {
	p.mu.Lock()
	defer p.mu.Unlock()
	if p.skip != nil {
		select {
		case p.skip <- struct{}{}:
		default:
		}
	}
}

func (p *Player) Stop() {
	p.mu.Lock()
	defer p.mu.Unlock()

	p.stopped = true

	if p.stop != nil {
		close(p.stop)
		p.stop = nil
		p.skip = nil
	}
	if p.current != nil {
		kill := exec.Command("taskkill", "/F", "/T", "/PID", fmt.Sprintf("%d", p.current.Pid))
		kill.Run()
		p.current = nil
	}
}

func Shuffle(items []emby.Item) []emby.Item {

	rand.Seed(time.Now().UnixNano())

	rand.Shuffle(len(items), func(i, j int) {
		items[i], items[j] = items[j], items[i]
	})

	return items
}
