# Emby-Lastfm-TUI
A simple Go program written by AI to play and scrobble music in your Emby server, with ~15MiB of RAM usage. 

Below only works on Windows; Linux users should be able to figure it out by themselves. (Mac? Really?)

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
go build ./cmd/emby-tui
``` 
If you didn't do anything wrong, a file called "emby-tui" should appear. Rename it to "emby-tui.exe". 
## Config 
Follow the instructions below or in the JSON file:
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
After you run the PowerShell script, everything should be filled, and you are good to go! (no pun intended)
## How to use 
- `Enter` or `Space` to play/select 
  - `n` for Normal playback 
  - `s` for Shuffle
- `s` is also for skipping songs
- `Backspace` to go back 
- `DEL` to clear the current queue
- `q` for exiting the program after the song finishes playing
- All others are written on there

## Known issues
- Sometimes when switching pages, song lists misalign and text overlaps; press `?` twice to refresh it.
- You can not select another song when another song is playing; press DEL before doing that
- mpv might use a lot of RAM; installing "Mem Reduct" helps.
- The shuffle mode sometimes shuffles too much, such is life, I suppose?
- Change the font of the command window if some words refuse to render.
### Changelog?
- ver2.0
    - Fixed many problems related to the display & progress bar
    - Added more functionality to the `DEL` and the `s` key 
- ver1.0 
    - ~~Where all miracles started, duh~~
## Other 
- Why upload a vibe-coded program? 
    - Others might need it; instead of vide-code another one, just use this one (works!), and save tokens. 
- There is an icon file in the folder, you can use it to change the boring(?) icon.
- You are free to contribute; either vibe-coded or skill-coded contributions work.
    - Be careful not to upload your api keys tho.