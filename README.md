# Emby-Lastfm-TUI
A simple Go program written by AI to play and scrobble music in your Emby server, with ~15MiB of RAM usage. 

Below only works on windows, Linux users should able to figure it out by yourself.

## Prerequisite
Go and mpv installed, and add to PATH 
(and git, of course)
### Verify
```bash 
go version 
```
and 
```bash 
mpv
```
## Compile
```bash 
git clone https://github.com/Farmertree8/Emby-Lastfm-TUI 
cd Emby-Lastfm-TUI
``` 
In the folder:
```bash 
go mod tidy 
go build -o emby-tui ./cmd/emby-tui
``` 
If you didn't do anything wrong, a file called "emby-tui" should appear. Rename it to "emby-tui.exe". 
## Config 
```json 
{
  "emby_url": "http://localhost:8096",
  "emby_api_key": "API-KEY-GENERATE-IN-CONTROL-PANEL",
  "user_id": "EDIT-USER-TO-SEE-IN-THE-URL",

  "lastfm_api_key": "APPLY-FROM-LASTFM",
  "lastfm_secret": "FILL-THE-SECRET-HERE",
  "lastfm_session_key": "RUN-THE-POWERSHELL-SCRIPT-AFTER-FINISHING-ABOVE-AND-SAVED"
}
``` 
After you run the PowerShell script, everything should be filled, and you are good to go (no pun intended)! 
## How to use 
- Enter to play/select 
  - n for Normal playback 
  - s for Shuffle
- Backspace to go back 
- Left click to pause?
- All others are written on there, duh

## Known issues
- The command line output overwrites the song list; press `?` twice to refresh it. 
- You can not select another song when another song is playing; just close & reopen the program. 
- mpv might use a lot of RAM; installing "Mem Reduct" helps.
- The shuffle mode sometimes shuffles too much, such is life, I suppose? 
## Other 
- Why upload a vibe-coded program? 
  Others might need it; instead of vide-code another one, just use this one (works!), and save tokens. 
- You are free to contribute; either vibe-coded or skill-coded contributions work.