package models

type RedditPost struct {
	Data struct {
		Title  string `json:"title"`
		Ups    int    `json:"ups"`
		Author string `json:"author"`
		ID     string `json:"id"`
	} `json:"data"`
}
