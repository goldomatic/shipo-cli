package twitter

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/dghubble/oauth1"
)

type TwitterClient struct {
	httpClient *http.Client
}

func NewTwitterClient(consumerKey, consumerSecret, accessToken, accessSecret string) *TwitterClient {
	config := oauth1.NewConfig(consumerKey, consumerSecret)
	token := oauth1.NewToken(accessToken, accessSecret)
	httpClient := config.Client(oauth1.NoContext, token)

	return &TwitterClient{httpClient: httpClient}
}

func (t *TwitterClient) CreateTweet(content string) error {
	url := "https://api.x.com/2/tweets" // Twitter API v2 endpoint
	payload := map[string]string{
		"text": content,
	}
	body, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to create payload: %w", err)
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(body))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := t.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("failed to post tweet, status: %s, body: %s", resp.Status, string(body))
	}

	fmt.Println("Tweet posted successfully!")
	return nil
}
