package stats

import (
	"sort"
	"vote-tracker/models"
)

type Statistics struct {
	MostUpvotedPost models.RedditPost
	UserPostCounts  map[string]int
}

type userPostCount struct {
	User  string
	Posts int
}

func NewStatistics() *Statistics {
	return &Statistics{
		UserPostCounts: make(map[string]int),
	}
}

func (s *Statistics) Update(post models.RedditPost) {
	// Update most upvoted post
	if post.Data.Ups > s.MostUpvotedPost.Data.Ups {
		s.MostUpvotedPost = post
	}

	// Update user post counts
	s.UserPostCounts[post.Data.Author]++
}

func (s *Statistics) GetTopUsers(n int) []userPostCount {
	var userPostCounts []userPostCount
	for user, count := range s.UserPostCounts {
		userPostCounts = append(userPostCounts, userPostCount{User: user, Posts: count})
	}

	// sort the slice by post count descending
	sort.Slice(userPostCounts, func(i, j int) bool {
		return userPostCounts[i].Posts > userPostCounts[j].Posts
	})

	if n > len(userPostCounts) {
		n = len(userPostCounts)
	}
	return userPostCounts[:n]
}
