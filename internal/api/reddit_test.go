package api

import (
	"net/http"
	"net/http/httptest"
	"os"
	"sync"
	"testing"
	"vote-tracker/models"

	"github.com/joho/godotenv"
)

func TestFetchPosts(t *testing.T) {
	if err := godotenv.Load("../../.env"); err != nil {
		t.Fatalf("Error loading .env file")
	}
	// Mock Reddit API server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{
			"data": {
				"after": "1",
				"children": [
					{
						"data": {
							"title": "Test Post",
							"ups": 100,
							"author": "test_user",
							"id": "1"
						}
					}
				]
			}
		}`))
	}))
	defer server.Close()

	// Override the Reddit API URL with the mock server URL
	os.Setenv("REDDIT_API_URL", server.URL)

	// Test fetchPosts
	postChan := make(chan models.RedditPost, 100) // Buffered channel to avoid blocking

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		err := FetchPosts("test", "mock_token", postChan)
		if err != nil {
			t.Errorf("fetchPosts failed: %v", err)
		}
	}()

	wg.Wait()
	close(postChan)

	// Now, consume from the channel and count the posts
	postCount := 0
	for post := range postChan {
		if post.Data.Title != "Test Post" {
			t.Fatalf("expected title 'Test Post', got '%s'", post.Data.Title)
		}
		postCount++
		t.Logf("Received post: %+v", post)
	}

	// Assert the number of posts
	expectedPostCount := 1 // Change this based on how many posts are sent
	if postCount != expectedPostCount {
		t.Errorf("Expected %d posts, but got %d", expectedPostCount, postCount)
	}

}

func TestGetOAuthToken(t *testing.T) {
	// load .env file in root directory

	if err := godotenv.Load("../../.env"); err != nil {
		t.Fatalf("Error loading .env file")
	}
	// Mock Reddit OAuth server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{
			"access_token": "mock_token"
		}`))
	}))
	defer server.Close()

	// Override the Reddit OAuth URL with the mock server URL
	os.Setenv("REDDIT_AUTH_URL", server.URL)

	// Test GetOAuthToken
	token, err := GetOAuthToken()
	if err != nil {
		t.Fatalf("GetOAuthToken failed: %v", err)
	}

	if token != "mock_token" {
		t.Fatalf("expected token 'mock_token', got '%s'", token)
	}
}
