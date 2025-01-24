package bluesky

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

type BlueskyClient struct {
	Jwt   string
	Did   string
	Actor string
}

func (b *BlueskyClient) Authenticate(handle, password string) error {
	authURL := "https://bsky.social/xrpc/com.atproto.server.createSession"
	authPayload := map[string]string{
		"identifier": handle,
		"password":   password,
	}

	payload, err := json.Marshal(authPayload)
	if err != nil {
		return fmt.Errorf("failed to create auth payload: %w", err)
	}

	req, err := http.NewRequest("POST", authURL, bytes.NewBuffer(payload))
	if err != nil {
		return fmt.Errorf("failed to create auth request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send auth request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("authentication failed: %s", body)
	}

	var authResponse AuthResponse
	if err := json.NewDecoder(resp.Body).Decode(&authResponse); err != nil {
		return fmt.Errorf("failed to decode auth response: %w", err)
	}

	b.Jwt = authResponse.AccessJwt
	b.Did = authResponse.Did
	b.Actor = handle
	return nil
}

func (b *BlueskyClient) CreatePost(content string) error {
	postURL := "https://bsky.social/xrpc/com.atproto.repo.createRecord"
	postRecord := PostRecord{
		Type:      "app.bsky.feed.post",
		Text:      content + " #Shipo-CLI",
		CreatedAt: time.Now().Format(time.RFC3339),
	}

	createRecordPayload := CreateRecordRequest{
		Repo:       b.Did,
		Collection: "app.bsky.feed.post",
		Record:     postRecord,
	}

	payload, err := json.Marshal(createRecordPayload)
	if err != nil {
		return fmt.Errorf("failed to create post payload: %w", err)
	}

	req, err := http.NewRequest("POST", postURL, bytes.NewBuffer(payload))
	if err != nil {
		return fmt.Errorf("failed to create post request: %w", err)
	}
	req.Header.Set("Authorization", "Bearer "+b.Jwt)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send post request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("post failed: %s", body)
	}

	fmt.Println("Post successful!")
	return nil
}

func (b *BlueskyClient) CountPostsToday() (int, error) {
	queryURL := fmt.Sprintf("https://bsky.social/xrpc/app.bsky.feed.getAuthorFeed?actor=%s", b.Actor)
	req, err := http.NewRequest("GET", queryURL, nil)
	if err != nil {
		return 0, fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Authorization", "Bearer "+b.Jwt)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return 0, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return 0, fmt.Errorf("failed to fetch posts: %s", body)
	}

	var response FeedResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return 0, fmt.Errorf("failed to decode response: %w", err)
	}

	count := 0
	today := time.Now().Format("2006-01-02")
	for _, post := range response.Feed {
		if strings.Contains(post.Record.Text, "#Shipo-CLI") && strings.HasPrefix(post.Record.CreatedAt, today) {
			count++
		}
	}

	return count, nil
}
