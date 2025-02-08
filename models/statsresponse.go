package models

// StatsResponse represents the response structure for statistics.
type StatsResponse struct {
	TopUsers []struct {
		User  string `json:"user"`
		Posts int    `json:"posts"`
	} `json:"top_users"`

	MostUpvotedPost struct {
		Data struct {
			Title  string `json:"title"`
			Ups    int    `json:"ups"`
			Author string `json:"author"`
			ID     string `json:"id"`
		} `json:"data"`
	} `json:"most_upvoted_post"`
}
