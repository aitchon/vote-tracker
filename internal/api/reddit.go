package api

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"
	"vote-tracker/internal/stats"
	"vote-tracker/models"
)

var lastFetchedPostID string

// GetOAuthToken fetches an OAuth2 token from Reddit
func GetOAuthToken() (string, error) {
	clientID := os.Getenv("REDDIT_CLIENT_ID")
	clientSecret := os.Getenv("REDDIT_CLIENT_SECRET")
	userAgent := os.Getenv("REDDIT_USER_AGENT")
	redditAuthURL := os.Getenv("REDDIT_AUTH_URL")
	if clientID == "" || clientSecret == "" {
		return "", fmt.Errorf("REDDIT_CLIENT_ID or REDDIT_CLIENT_SECRET not set")
	}

	client := &http.Client{}
	data := url.Values{}
	data.Set("grant_type", "client_credentials")

	req, err := http.NewRequest("POST", redditAuthURL, strings.NewReader(data.Encode()))
	if err != nil {
		return "", err
	}

	req.SetBasicAuth(clientID, clientSecret)
	req.Header.Set("User-Agent", userAgent)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("failed to get OAuth token: %s", string(body))
	}

	var result struct {
		AccessToken string `json:"access_token"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", err
	}

	return result.AccessToken, nil
}

// FetchPosts fetches posts from a subreddit using OAuth2 token
func FetchPosts(subreddit, token string, postChan chan<- models.RedditPost) error {
	userAgent := os.Getenv("REDDIT_USER_AGENT")
	redditPostsURL := os.Getenv("REDDIT_API_URL")
	var url string
	// Check if the URL contains a placeholder for the subreddit
	if strings.Contains(redditPostsURL, "%s") {
		url = fmt.Sprintf(redditPostsURL, subreddit)
		if lastFetchedPostID != "" {
			url += "?after=" + lastFetchedPostID // Use the 'after' parameter to fetch newer posts
		}
	} else {
		// if no placeholder, use the URL as is
		url = redditPostsURL
	}

	// Retry mechanism
	maxRetries := 3
	for retry := 0; retry < maxRetries; retry++ {
		req, err := http.NewRequest("GET", url, nil)
		if err != nil {
			return err
		}

		req.Header.Set("User-Agent", userAgent)
		req.Header.Set("Authorization", "Bearer "+token)

		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			return err
		}
		defer resp.Body.Close()

		// Handle rate limiting (429 Too Many Requests)
		if resp.StatusCode == http.StatusTooManyRequests {
			// Parse rate limit headers
			rateLimitReset := resp.Header.Get("X-Ratelimit-Reset")
			rateLimitUsed := resp.Header.Get("X-Ratelimit-Used")
			rateLimitRemaining := resp.Header.Get("X-Ratelimit-Remaining")
			resetSeconds, err := strconv.Atoi(rateLimitReset)
			if err != nil {
				log.Printf("Failed to parse X-Ratelimit-Reset header: %v", err)
				resetSeconds = 60 // Default to 60 seconds if parsing fails
			}

			// Wait for the rate limit to reset
			waitTime := time.Duration(resetSeconds) * time.Second

			log.Printf("Rate Limit exceeded: Retrying after %v...Used=%s, Remaining=%s, Reset=%s",
				waitTime, rateLimitUsed, rateLimitRemaining, rateLimitReset)

			time.Sleep(waitTime)

			continue // Retry the request
		}

		// Handle other non-200 status codes
		if resp.StatusCode != http.StatusOK {
			return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
		}

		// Parse the response body
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return err
		}

		var redditResponse struct {
			Data struct {
				After    string              `json:"after"`
				Children []models.RedditPost `json:"children"`
			} `json:"data"`
		}
		if err := json.Unmarshal(body, &redditResponse); err != nil {
			return err
		}

		// Send posts to the channel
		for _, post := range redditResponse.Data.Children {
			postChan <- post
		}

		// Update the last fetched post ID
		if len(redditResponse.Data.Children) > 0 {
			lastFetchedPostID = redditResponse.Data.After
		}

		return nil // Success
	}

	// If all retries fail, return an error
	return fmt.Errorf("failed to fetch posts after %d retries", maxRetries)
}

// ProcessPosts processes posts from the channel and updates the statistics
func ProcessPosts(postChan <-chan models.RedditPost, stats *stats.Statistics) {
	for post := range postChan {
		stats.Update(post)
		log.Printf("Processed post: %s (%d upvotes)", post.Data.Title, post.Data.Ups)
	}
}
