package slackhook

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"

	"github.com/nlopes/slack"
)

// Client is used for making HTTP requests to a Slack Incoming Webhook.
type Client struct {
	HTTPClient *http.Client
	url        string
}

// NewClient returns a new Client.
func NewClient(webhookURL string) (*Client, error) {
	u, err := url.Parse(webhookURL)
	if err != nil {
		return nil, err
	}

	return &Client{
		HTTPClient: &http.Client{
			Timeout: time.Second * 5,
		},
		url: u.String(),
	}, nil
}

// SendMessage sends a Slack message to an Incoming Webhook.
func (c *Client) SendMessage(msg slack.Msg) error {
	buf := &bytes.Buffer{}
	enc := json.NewEncoder(buf)

	if err := enc.Encode(msg); err != nil {
		return err
	}

	resp, err := c.HTTPClient.Post(c.url, "application/json", buf)
	if err != nil {
		return fmt.Errorf("slackhook: could not send HTTP request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		return nil
	}

	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("slackhook: could not read HTTP response body: %v",
			err,
		)
	}

	return fmt.Errorf(
		"slackhook: received erroneous response (code: %v, body: %v)",
		resp.StatusCode,
		respBody,
	)
}
