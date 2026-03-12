package lastfm

import (
 "crypto/md5"
 "encoding/hex"
 "fmt"
 "net/http"
 "net/url"
 "sort"
)

type Client struct {
 ApiKey string
 Secret string
 SessionKey string
 Client *http.Client
}

func New(key, secret, session string) *Client {

 return &Client{
  ApiKey: key,
  Secret: secret,
  SessionKey: session,
  Client: &http.Client{},
 }
}

func (c *Client) sign(params url.Values) string {
 
 keys := make([]string, 0, len(params))
 for k := range params {
     keys = append(keys, k)
 }
 sort.Strings(keys)
 
 raw := ""
 for _, k := range keys {
     raw += k + params[k][0]
 }
 raw += c.Secret

 hash := md5.Sum([]byte(raw))

 return hex.EncodeToString(hash[:])
}

func (c *Client) post(params url.Values) error {

 params.Set("api_key", c.ApiKey)
 params.Set("sk", c.SessionKey)

 sig := c.sign(params)

 params.Set("api_sig", sig)

 _, err := c.Client.PostForm(
  "https://ws.audioscrobbler.com/2.0/",
  params,
 )

 return err
}

func (c *Client) UpdateNowPlaying(track, artist string) error {

 params := url.Values{
  "method": {"track.updateNowPlaying"},
  "track": {track},
  "artist": {artist},
 }

 return c.post(params)
}

func (c *Client) Scrobble(track, artist string, ts int64) error {

 params := url.Values{
  "method": {"track.scrobble"},
  "track": {track},
  "artist": {artist},
  "timestamp": {fmt.Sprint(ts)},
 }

 return c.post(params)
}