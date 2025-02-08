package stats

import "vote-tracker/models"

type Statistics struct {
	MostUpvotedPost models.RedditPost
	UserPostCounts  map[string]int
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

func (s *Statistics) GetTopUser() string {
	topUser := ""
	maxPosts := 0
	for user, count := range s.UserPostCounts {
		if count > maxPosts {
			topUser = user
			maxPosts = count
		}
	}
	return topUser
}

func (s *Statistics) GetTopUserCount() int {
	maxPosts := 0
	for _, count := range s.UserPostCounts {
		if count > maxPosts {
			maxPosts = count
		}
	}
	return maxPosts
}
