package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

const hnAPIBase = "https://hacker-news.firebaseio.com/v0"

// HackerNewsClient interacts with the Hacker News Firebase API.
type HackerNewsClient struct {
	httpClient *http.Client
}

// NewHackerNewsClient creates a new HN API client.
func NewHackerNewsClient() *HackerNewsClient {
	return &HackerNewsClient{
		httpClient: &http.Client{},
	}
}

// Story represents a Hacker News story/item.
type Story struct {
	ID          int    `json:"id"`
	Title       string `json:"title"`
	URL         string `json:"url,omitempty"`
	Score       int    `json:"score"`
	Author      string `json:"by"`
	Time        int64  `json:"time"`
	Text        string `json:"text,omitempty"`
	Descendants int    `json:"descendants"`
	Type        string `json:"type"`
}

// GetTopStories fetches the top story IDs and returns the first `limit` stories.
func (c *HackerNewsClient) GetTopStories(ctx context.Context, limit int) ([]Story, error) {
	// Fetch top story IDs
	ids, err := c.getTopStoryIDs(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get top story IDs: %w", err)
	}

	if limit > len(ids) {
		limit = len(ids)
	}

	// Fetch individual stories
	stories := make([]Story, 0, limit)
	for i := 0; i < limit; i++ {
		story, err := c.GetStory(ctx, ids[i])
		if err != nil {
			// Skip failed fetches, continue with others
			continue
		}
		stories = append(stories, *story)
	}

	return stories, nil
}

// GetStory fetches a specific story by ID.
func (c *HackerNewsClient) GetStory(ctx context.Context, id int) (*Story, error) {
	url := fmt.Sprintf("%s/item/%d.json", hnAPIBase, id)

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch story: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API returned status %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	var story Story
	if err := json.Unmarshal(body, &story); err != nil {
		return nil, fmt.Errorf("failed to parse story: %w", err)
	}

	return &story, nil
}

// getTopStoryIDs fetches the list of top story IDs from HN.
func (c *HackerNewsClient) getTopStoryIDs(ctx context.Context) ([]int, error) {
	url := fmt.Sprintf("%s/topstories.json", hnAPIBase)

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch top stories: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API returned status %d", resp.StatusCode)
	}

	var ids []int
	if err := json.NewDecoder(resp.Body).Decode(&ids); err != nil {
		return nil, fmt.Errorf("failed to parse IDs: %w", err)
	}

	return ids, nil
}
