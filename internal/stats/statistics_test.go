package stats

import (
	"testing"
	"vote-tracker/models"
)

func TestStatistics(t *testing.T) {
	stats := NewStatistics()

	// Test Update
	post1 := models.RedditPost{
		Data: struct {
			Title  string `json:"title"`
			Ups    int    `json:"ups"`
			Author string `json:"author"`
			ID     string `json:"id"`
		}{
			Title:  "Post 1",
			Ups:    100,
			Author: "user1",
			ID:     "1",
		},
	}
	post2 := models.RedditPost{
		Data: struct {
			Title  string `json:"title"`
			Ups    int    `json:"ups"`
			Author string `json:"author"`
			ID     string `json:"id"`
		}{
			Title:  "Post 2",
			Ups:    200,
			Author: "user2",
			ID:     "2",
		},
	}
	post3 := models.RedditPost{
		Data: struct {
			Title  string `json:"title"`
			Ups    int    `json:"ups"`
			Author string `json:"author"`
			ID     string `json:"id"`
		}{
			Title:  "Post 3",
			Ups:    150,
			Author: "user1",
			ID:     "3",
		},
	}

	stats.Update(post1)
	stats.Update(post2)
	stats.Update(post3)

	// Test MostUpvotedPost
	if stats.MostUpvotedPost.Data.Title != "Post 2" {
		t.Fatalf("expected most upvoted post 'Post 2', got '%s'", stats.MostUpvotedPost.Data.Title)
	}

	// Test GetTopUser
	if stats.GetTopUser() != "user1" {
		t.Fatalf("expected top user 'user1', got '%s'", stats.GetTopUser())
	}

	// Test GetTopUserCount
	if stats.GetTopUserCount() != 2 {
		t.Fatalf("expected top user count 2, got %d", stats.GetTopUserCount())
	}
}
