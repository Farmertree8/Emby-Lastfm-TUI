package emby

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type Client struct {
	BaseURL string
	ApiKey  string
	UserID  string
	Client  *http.Client
}

func New(url string, key string, uid string) *Client {
	return &Client{
		BaseURL: url,
		ApiKey:  key,
		UserID:  uid,
		Client:  &http.Client{},
	}
}

func (c *Client) GetLibraries() ([]Item, error) {

	url := fmt.Sprintf(
		"%s/Users/%s/Views?api_key=%s",
		c.BaseURL,
		c.UserID,
		c.ApiKey,
	)

	resp, err := c.Client.Get(url)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	var result ItemResponse

	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		return nil, err
	}

	return result.Items, nil
}

func (c *Client) GetItems(parent string) ([]Item, error) {

	url := fmt.Sprintf(
		"%s/Users/%s/Items?ParentId=%s&api_key=%s&SortBy=DateCreated&SortOrder=Descending",
		c.BaseURL,
		c.UserID,
		parent,
		c.ApiKey,
	)

	resp, err := c.Client.Get(url)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	var result ItemResponse

	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		return nil, err
	}

	return result.Items, nil
}

func (c *Client) StreamURL(id string) string {

	return fmt.Sprintf(
		"%s/Audio/%s/stream?static=true&api_key=%s",
		c.BaseURL,
		id,
		c.ApiKey,
	)
}

func (c *Client) GetSessions() ([]Session, error) {

	url := fmt.Sprintf(
		"%s/Sessions?api_key=%s",
		c.BaseURL,
		c.ApiKey,
	)

	resp, err := c.Client.Get(url)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	var sessions []Session

	err = json.NewDecoder(resp.Body).Decode(&sessions)
	if err != nil {
		return nil, err
	}

	return sessions, nil
}