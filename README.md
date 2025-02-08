# Reddit Vote Tracker Application

This application fetches posts from Reddit, tracks statistics (e.g., most upvoted posts, active users), and exposes the data via a REST API. It is written in Go.

---

## Table of Contents

1. [Prerequisites](#prerequisites)
2. [Setup](#setup)
3. [Configuration](#configuration)
4. [Running the Application](#running-the-application)
5. [REST API Endpoints](#rest-api-endpoints)
6. [Graceful Shutdown](#graceful-shutdown)

---

## Prerequisites

Before running the application, ensure you have the following installed:

- **Go** (version 1.20 or higher): [Install Go](https://golang.org/doc/install)
- **Reddit API Credentials**:
  - Create a Reddit app at [Reddit App Preferences](https://www.reddit.com/prefs/apps).
  - Obtain the `Client ID` and `Client Secret`.

---

## Setup

Clone the repository:
   ```bash
   git clone https://github.com/aitchon/vote-tracker.git
   cd vote-tracker
Install dependencies:

```bash
go mod tidy
```
Create a .env file in the root directory and add your Reddit Client ID and Client Secret credentials. Use the values for REDDIT_API_URL, REDDIT_AUTH_URL, REDDIT_USER_AGENT here:

```bash
REDDIT_CLIENT_ID=your_client_id
REDDIT_CLIENT_SECRET=your_client_secret
REDDIT_API_URL=https://oauth.reddit.com/r/%s/new.json
REDDIT_AUTH_URL=https://www.reddit.com/api/v1/access_token
REDDIT_USER_AGENT=RedditStatsApp/0.1
SUBREDDIT_NAME=golang
```
# Configuration
The application can be configured using environment variables. The following variables are required:

| Variable | Description |
|----------|----------|
| REDDIT_CLIENT_ID | Your Reddit app's Client ID |
| REDDIT_CLIENT_SECRET | Your Reddit app's Client Secret |
| REDDIT_API_URL | Reddit posts URL |
| REDDIT_AUTH_URL | Reddit authentication URL |
| REDDIT_USER_AGENT | User agent to use when calling Reddit endpoints |
| SUBREDDIT_NAME | The subreddit to fetch posts from |

# Running the Application
Build the application:

```bash
go build -o vote-tracker main.go
```
Run the application:

```bash
./vote-tracker
```
Alternatively, you can run it directly without building:

```bash
go run main.go
```
The application will:

* Fetch posts from the specified subreddit.
* Track statistics in memory.
* Start a REST API server on http://localhost:8080/.

# REST API Endpoint
The application exposes the following REST API endpoint:

Get Statistics
Endpoint: GET /stats

Description: Returns the most upvoted post, top user, and total posts for top user.

Example Response:

```json
{
  "most_upvoted_post": {
    "title": "Test Post",
    "ups": 100,
    "author": "test_user"
  },
  "top_user": "test_user",
  "top_user_posts": 5
}
```

# Running Tests
To run the tests, use the following command:

```bash
go test ./...
```
This will run the unit tests for the application.

# Graceful Shutdown
The application supports graceful shutdown. When you press Ctrl+C, it will:

* Stop fetching new posts.
* Finish processing existing posts.
* Shut down the REST API server.

