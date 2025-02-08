package main

import (
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
	"vote-tracker/internal/api"
	"vote-tracker/internal/server"
	"vote-tracker/internal/stats"
	"vote-tracker/models"

	"github.com/joho/godotenv"
)

func main() {
	// Load REDDIT_CLIENT_ID, REDDIT_CLIENT_SECRET, REDDIT_API_URL,
	// REDDIT_AUTH_URL, REDDIT_USER_AGENT from .env file
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}

	token, err := api.GetOAuthToken()
	if err != nil {
		log.Fatalf("Failed to get OAuth token: %v", err)
	}

	statsTracker := stats.NewStatistics()

	// Create a channel for posts
	postChan := make(chan models.RedditPost, 100) // Buffered channel to avoid blocking

	doneChan := make(chan struct{}) // Channel to signal goroutines to stop

	// Create a WaitGroup to wait for goroutines to finish
	var wg sync.WaitGroup

	// Start the post processor
	wg.Add(1)
	go func() {
		defer wg.Done()
		api.ProcessPosts(postChan, statsTracker)
	}()

	// Start a goroutine to fetch posts periodically
	subreddit := os.Getenv("SUBREDDIT_NAME")
	wg.Add(1)
	go func() {
		defer wg.Done()
		for {
			select {
			case <-doneChan: // Check for shutdown signal
				log.Println("Stopping post fetcher...")
				return
			default:
				if err := api.FetchPosts(subreddit, token, postChan); err != nil {
					log.Printf("Error fetching posts: %v", err)
				}
				time.Sleep(30 * time.Second)
			}
		}
	}()

	// Start the REST API server
	wg.Add(1)
	go func() {
		defer wg.Done()
		port := "8080"
		log.Printf("Starting API server on port %s...", port)
		log.Fatal(http.ListenAndServe(":"+port, server.StartAPIServer(statsTracker)))
		server.StartAPIServer(statsTracker)
	}()

	// Listen for Ctrl+C (SIGINT)
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	// Wait for interrupt signal
	<-sigChan
	log.Println("Received interrupt signal. Shutting down...")

	// Close the channel to signal goroutines to stop
	close(postChan)

	// Wait for all goroutines to finish
	wg.Wait()

	log.Println("Application shutdown complete.")
}
