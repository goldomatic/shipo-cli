package main

import (
	"flag"
	"fmt"
	"os"
	"shipo-cli/bluesky"
	"shipo-cli/twitter"
	"shipo-cli/utils"
	"strconv"
	"strings"
)

func main() {
	content := flag.String("c", "", "Content of the post (required)")
	platform := flag.String("p", "b", "Platform to post on: 'b' for Bluesky, 't' for Twitter, 'bt' or 'tb' for both (default: 'b')")
	flag.Parse()

	if *content == "" {
		fmt.Println("Error: -c is required.")
		flag.Usage()
		os.Exit(1)
	}

	// Load configuration
	config, err := utils.LoadConfig()
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}

	// Handle posting to Bluesky
	if strings.Contains(*platform, "b") {
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

		client := bluesky.BlueskyClient{}
		fmt.Println("Authenticating with Bluesky...")
		if err := client.Authenticate(handle, password); err != nil {
			fmt.Printf("Error: %v\n", err)
			os.Exit(1)
		}

		postCount, err := client.CountPostsToday()
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			os.Exit(1)
		}

		if limit > 0 && postCount >= limit {
			fmt.Printf("Daily post limit of %d reached. Try again tomorrow!\n", limit)
			os.Exit(1)
		}

		fmt.Println("Creating post on Bluesky...")
		if err := client.CreatePost(*content); err != nil {
			fmt.Printf("Error: %v\n", err)
			os.Exit(1)
		}
	}

	if strings.Contains(*platform, "t") {
		consumerKey := config["twitter_consumer_key"]
		consumerSecret := config["twitter_consumer_secret"]
		accessToken := config["twitter_access_token"]
		accessSecret := config["twitter_access_secret"]

		if consumerKey == "" || consumerSecret == "" || accessToken == "" || accessSecret == "" {
			fmt.Println("Error: Missing Twitter API credentials in the config file")
			os.Exit(1)
		}

		twitterClient := twitter.NewTwitterClient(consumerKey, consumerSecret, accessToken, accessSecret)
		fmt.Println("Creating post on Twitter...")
		if err := twitterClient.CreateTweet(*content); err != nil {
			fmt.Printf("Error: %v\n", err)
			os.Exit(1)
		}
	}
}
