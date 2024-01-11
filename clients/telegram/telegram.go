package telegram

import (
	"encoding/json"
	"io"
	"links_tg-bot/lib/e"
	"net/http"
	"net/url"
	"path"
	"strconv"
)

// type Client is used for connection with API service telegram

type Client struct {
	host     string // host is host of API service telegram (on example below tg-bot.com is host)
	basePath string // basePath is prefix which starts all requests (example: tg-bot.com/bot<token>; bot<token> is basePath)
	client   http.Client
}

const (
	getUpdatesMethod  = "getUpdates"
	sendMessageMethod = "sendMessage"
)

func NewClient(host string, token string) *Client {
	return &Client{
		host:     host,
		basePath: newBasePath(token),
		client:   http.Client{},
	}
}

func newBasePath(token string) string {
	return "bot" + token
}

func (c *Client) Updates(offset, limit int) ([]Update, error) {

	q := url.Values{}

	q.Add("offset", strconv.Itoa(offset))
	q.Add("limit", strconv.Itoa(limit))

	// do request <-getUpdates
	data, err := c.doRequest(getUpdatesMethod, q)

	if err != nil {
		return nil, err
	}

	var res UpdatesResponse

	if err = json.Unmarshal(data, &res); err != nil {
		return nil, err
	}
	return res.Result, nil

}

func (c *Client) SendMessages(chatID int, text string) error {
	q := url.Values{}

	q.Add("chat_id", strconv.Itoa(chatID))
	q.Add("text", text)

	_, err := c.doRequest(sendMessageMethod, q)

	if err != nil {
		return e.Wrap("can't send request", err)
	}
	return nil

}

func (c *Client) doRequest(method string, query url.Values) (data []byte, err error) {
	defer func() {
		err = e.WrapIfErr("can't do request", err)
	}()

	u := url.URL{
		Scheme: "https",
		Host:   c.host,
		Path:   path.Join(c.basePath, method),
	}

	req, err := http.NewRequest(http.MethodGet, u.String(), nil)

	if err != nil {
		return nil, err
	}

	req.URL.RawQuery = query.Encode()

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}

	defer func() { _ = resp.Body.Close() }()

	data, err = io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return
}
