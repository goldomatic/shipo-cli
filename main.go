package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/user"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

type AuthResponse struct {
	AccessJwt string `json:"accessJwt"`
	Did       string `json:"did"`
}

type PostRecord struct {
	Type      string `json:"$type"`
	Text      string `json:"text"`
	CreatedAt string `json:"createdAt"`
}

type CreateRecordRequest struct {
	Repo       string     `json:"repo"`
	Collection string     `json:"collection"`
	Record     PostRecord `json:"record"`
}

// getConfigFilePath returns the path to the config file.
func getConfigFilePath() (string, error) {
	usr, err := user.Current()
	if err != nil {
		return "", fmt.Errorf("failed to get current user: %w", err)
	}
	return filepath.Join(usr.HomeDir, ".config", "shipo", "config"), nil
}

// loadConfig loads the configuration file, ignores comments, and trims spaces.
func loadConfig() (map[string]string, error) {
	configPath, err := getConfigFilePath()
	if err != nil {
		return nil, err
	}

	file, err := os.Open(configPath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, errors.New("config file not found; please create ~/.config/shipo/config with your handle, password, and limit")
		}
		return nil, fmt.Errorf("failed to open config file: %w", err)
	}
	defer file.Close()

	config := make(map[string]string)
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			return nil, fmt.Errorf("invalid config line: %s", line)
		}
		key := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])
		config[key] = value
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	return config, nil
}

// countPostsToday checks how many posts with #Shipo-CLI were made today.
func countPostsToday(jwt, did string, handle string) (int, error) {
	queryURL := fmt.Sprintf("https://bsky.social/xrpc/app.bsky.feed.getAuthorFeed?actor=%s", handle)

	req, err := http.NewRequest("GET", queryURL, nil)
	if err != nil {
		return 0, fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Authorization", "Bearer "+jwt)

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

	var response struct {
		Feed []struct {
			Record struct {
				Text      string `json:"text"`
				CreatedAt string `json:"createdAt"`
			} `json:"record"`
		} `json:"feed"`
	}
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

// authenticate retrieves an access token and DID from Bluesky.
func authenticate(handle, password string) (string, string, error) {
	authURL := "https://bsky.social/xrpc/com.atproto.server.createSession"

	authPayload := map[string]string{
		"identifier": handle,
		"password":   password,
	}
	payload, err := json.Marshal(authPayload)
	if err != nil {
		return "", "", fmt.Errorf("failed to create auth payload: %w", err)
	}

	req, err := http.NewRequest("POST", authURL, bytes.NewBuffer(payload))
	if err != nil {
		return "", "", fmt.Errorf("failed to create auth request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", "", fmt.Errorf("failed to send auth request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return "", "", fmt.Errorf("authentication failed: %s", body)
	}

	var authResponse AuthResponse
	if err := json.NewDecoder(resp.Body).Decode(&authResponse); err != nil {
		return "", "", fmt.Errorf("failed to decode auth response: %w", err)
	}

	return authResponse.AccessJwt, authResponse.Did, nil
}

// createPost sends a post to Bluesky.
func createPost(jwt, did, content string) error {
	postURL := "https://bsky.social/xrpc/com.atproto.repo.createRecord"

	postRecord := PostRecord{
		Type:      "app.bsky.feed.post",
		Text:      content + " #Shipo-CLI", // Add default tag
		CreatedAt: time.Now().Format(time.RFC3339),
	}

	createRecordPayload := CreateRecordRequest{
		Repo:       did,
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
	req.Header.Set("Authorization", "Bearer "+jwt)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send post request: %w", err)
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("post failed: %s", body)
	}

	fmt.Println("Post successful!")
	return nil
}

func main() {
	content := flag.String("content", "", "Content of the post (required)")
	flag.Parse()

	if *content == "" {
		fmt.Println("Error: --content is required.")
		flag.Usage()
		os.Exit(1)
	}

	config, err := loadConfig()
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}

	handle, ok := config["handle"]
	if !ok {
		fmt.Println("Error: 'handle' not found in the config file")
		os.Exit(1)
	}

	password, ok := config["password"]
	if !ok {
		fmt.Println("Error: 'password' not found in the config file")
		os.Exit(1)
	}

	limitStr, ok := config["limit"]
	if !ok {
		fmt.Println("Error: 'limit' not found in the config file")
		os.Exit(1)
	}

	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit < 0 {
		fmt.Println("Error: 'limit' must be a non-negative integer")
		os.Exit(1)
	}

	fmt.Println("Authenticating...")
	token, did, err := authenticate(handle, password)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}

	postCount, err := countPostsToday(token, did, handle)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}

	if limit > 0 && postCount >= limit {
		fmt.Printf("Daily post limit of %d reached. Try again tomorrow!\n", limit)
		os.Exit(1)
	}

	fmt.Println("Creating post...")
	if err := createPost(token, did, *content); err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}
}
