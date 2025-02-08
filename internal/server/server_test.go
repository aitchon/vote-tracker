package server

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"vote-tracker/internal/stats"
	"vote-tracker/models"
)

func TestAPIServer(t *testing.T) {
	// Create a mock Statistics instance
	statsTracker := stats.NewStatistics()
	statsTracker.Update(models.RedditPost{
		Data: struct {
			Title  string `json:"title"`
			Ups    int    `json:"ups"`
			Author string `json:"author"`
			ID     string `json:"id"`
		}{
			Title:  "Test Post",
			Ups:    100,
			Author: "test_user",
			ID:     "1",
		},
	})

	// Start the test server
	server := httptest.NewServer(StartAPIServer(statsTracker))
	defer server.Close()

	// Test the /stats endpoint
	resp, err := http.Get(server.URL + "/stats")
	if err != nil {
		t.Fatalf("failed to make request: %v", err)
	}
	defer resp.Body.Close()

	var statsResponse models.StatsResponse
	if err := readJSON(resp.Body, &statsResponse); err != nil {
		t.Fatalf("failed to read response: %v", err)
	}
	if statsResponse.MostUpvotedPost.Data.Ups != 100 {
		t.Fatalf("expected total votes to be 100, got %d", statsResponse.MostUpvotedPost.Data.Ups)
	}
	if statsResponse.MostUpvotedPost.Data.Title != "Test Post" {
		t.Fatalf("expected top post title  to \"Test Post\", got %s", statsResponse.MostUpvotedPost.Data.Title)
	}
	if statsResponse.MostUpvotedPost.Data.Author != "test_user" {
		t.Fatalf("expected top post author to be \"test_user\", got %s", statsResponse.MostUpvotedPost.Data.Author)
	}
	if statsResponse.TopUsers[0].User != "test_user" {
		t.Fatalf("expected top user to be \"test_user\", got %s", statsResponse.TopUsers[0].User)
	}
	if statsResponse.TopUsers[0].Posts != 1 {
		t.Fatalf("expected top user post count to be 1, got %d", statsResponse.TopUsers[0].Posts)
	}

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected status code 200, got %d", resp.StatusCode)
	}
}

// readJSON reads JSON from an io.Reader and decodes it into the provided interface.
func readJSON(r io.Reader, v interface{}) error {
	return json.NewDecoder(r).Decode(v)
}
