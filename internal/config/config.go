package config

import (
 "encoding/json"
 "os"
)

type Config struct {
 EmbyURL string `json:"emby_url"`
 EmbyAPIKey string `json:"emby_api_key"`
 UserID string `json:"user_id"`

 LastFMApiKey string `json:"lastfm_api_key"`
 LastFMSecret string `json:"lastfm_secret"`
 LastFMSessionKey string `json:"lastfm_session_key"`
}

func Load(path string) (*Config, error) {

 data, err := os.ReadFile(path)
 if err != nil {
  return nil, err
 }

 var cfg Config

 err = json.Unmarshal(data, &cfg)
 if err != nil {
  return nil, err
 }

 return &cfg, nil
}