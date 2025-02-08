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
	post4 := models.RedditPost{
		Data: struct {
			Title  string `json:"title"`
			Ups    int    `json:"ups"`
			Author string `json:"author"`
			ID     string `json:"id"`
		}{
			Title:  "Post 4",
			Ups:    50,
			Author: "user2",
			ID:     "4",
		},
	}
	post5 := models.RedditPost{
		Data: struct {
			Title  string `json:"title"`
			Ups    int    `json:"ups"`
			Author string `json:"author"`
			ID     string `json:"id"`
		}{
			Title:  "Post 5",
			Ups:    25,
			Author: "user2",
			ID:     "5",
		},
	}

	stats.Update(post1)
	stats.Update(post2)
	stats.Update(post3)
	stats.Update(post4)
	stats.Update(post5)

	// Test MostUpvotedPost
	if stats.MostUpvotedPost.Data.Title != "Post 2" {
		t.Fatalf("expected most upvoted post 'Post 2', got '%s'", stats.MostUpvotedPost.Data.Title)
	}

	// Test GetTopUser
	// write a unit test for GetTopUsers
	topUser := stats.GetTopUsers(10)[0].User
	if topUser != "user2" {
		t.Fatalf("expected top user 'user2', got '%s'", topUser)
	}

	// Test GetTopUserCount
	if stats.GetTopUserCount() != 3 {
		t.Fatalf("expected top user count 3, got %d", stats.GetTopUserCount())
	}
}
