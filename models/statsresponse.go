package models

// StatsResponse represents the response structure for statistics.
type StatsResponse struct {
	TopUserPosts    int    `json:"top_user_posts"`
	TopUser         string `json:"top_user"`
	MostUpvotedPost struct {
		Data struct {
			Title  string `json:"title"`
			Ups    int    `json:"ups"`
			Author string `json:"author"`
			ID     string `json:"id"`
		} `json:"data"`
	} `json:"most_upvoted_post"`
}
