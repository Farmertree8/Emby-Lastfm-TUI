package ui

import "emby-tui-scrobbler/internal/emby"

type listItem emby.Item

func (i listItem) Title() string {

 switch i.Type {

 case "Folder":
  return "Folder  " + i.Name

 case "MusicAlbum":
  return "Album   " + i.Name

 case "Audio":
  return "Song    " + i.Name

 default:
  return i.Name
 }
}

func (i listItem) Description() string {
 return i.Type
}

func (i listItem) FilterValue() string {
 return i.Name
}