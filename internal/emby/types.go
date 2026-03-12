package emby

type ItemResponse struct {
 Items []Item `json:"Items"`
}

type Item struct {
 Id string `json:"Id"`
 Name string `json:"Name"`
 Type string `json:"Type"`

 RunTimeTicks int64 `json:"RunTimeTicks"`

 Artists []string `json:"Artists"`
}

type Session struct {
 NowPlayingItem Item `json:"NowPlayingItem"`
 PlayState PlayState `json:"PlayState"`
}

type PlayState struct {
 PositionTicks int64 `json:"PositionTicks"`
}