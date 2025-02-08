package server

import (
	"encoding/json"
	"net/http"
	"vote-tracker/internal/stats"
)

func StartAPIServer(stats *stats.Statistics) http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("/stats", func(w http.ResponseWriter, r *http.Request) {
		response := map[string]interface{}{
			"most_upvoted_post": stats.MostUpvotedPost,
			"top_user":          stats.GetTopUser(),
			"top_user_posts":    stats.GetTopUserCount(),
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	})

	return mux
}
