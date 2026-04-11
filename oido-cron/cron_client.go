package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
)

const defaultAPIBase = "http://localhost:8080"

// CronClient communicates with the oido-studio API for cron management.
type CronClient struct {
	apiBase string
	token   string
	client  *http.Client
}

func NewCronClient() *CronClient {
	apiBase := os.Getenv("OIDO_API_BASE")
	if apiBase == "" {
		apiBase = defaultAPIBase
	}
	token := os.Getenv("OIDO_API_TOKEN")

	return &CronClient{
		apiBase: apiBase,
		token:   token,
		client:  &http.Client{},
	}
}

func (c *CronClient) doRequest(method, path string, body interface{}) ([]byte, error) {
	var reqBody io.Reader
	if body != nil {
		data, _ := json.Marshal(body)
		reqBody = bytes.NewReader(data)
	}

	req, err := http.NewRequest(method, c.apiBase+path, reqBody)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	if c.token != "" {
		req.Header.Set("Authorization", "Bearer "+c.token)
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	if resp.StatusCode >= 400 {
		return nil, fmt.Errorf("API error %d: %s", resp.StatusCode, string(respBody))
	}

	return respBody, nil
}

// CronJob represents a cron job
type CronJob struct {
	ID          int64  `json:"id"`
	Name        string `json:"name"`
	Schedule    string `json:"schedule"`
	Enabled     bool   `json:"enabled"`
	Delivery    string `json:"delivery"`
	Payload     string `json:"payload"`
}

// CronRunLog represents execution log
type CronRunLog struct {
	ID          int64  `json:"id"`
	JobID       int64  `json:"jobId"`
	StartedAt   string `json:"startedAt"`
	CompletedAt string `json:"completedAt,omitempty"`
	Status      string `json:"status"`
	Output      string `json:"output,omitempty"`
	Error       string `json:"error,omitempty"`
}

// ListJobs returns all cron jobs
func (c *CronClient) ListJobs() ([]CronJob, error) {
	resp, err := c.doRequest("GET", "/api/cron-jobs", nil)
	if err != nil {
		return nil, err
	}

	var jobs []CronJob
	if err := json.Unmarshal(resp, &jobs); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}
	return jobs, nil
}

// GetJob returns a specific job by ID
func (c *CronClient) GetJob(id int64) (*CronJob, error) {
	resp, err := c.doRequest("GET", fmt.Sprintf("/api/cron-jobs/%d", id), nil)
	if err != nil {
		return nil, err
	}

	var job CronJob
	if err := json.Unmarshal(resp, &job); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}
	return &job, nil
}

// CreateJobRequest is the request body for creating a job
type CreateJobRequest struct {
	Name     string `json:"name"`
	Schedule string `json:"schedule"`
	TZ       string `json:"tz,omitempty"`
	Message  string `json:"message"`
	Model    string `json:"model,omitempty"`
	Session  string `json:"session,omitempty"`
	Delivery string `json:"delivery,omitempty"`
	Channel  string `json:"channel,omitempty"`
	To       string `json:"to,omitempty"`
}

// CreateJob creates a new cron job
func (c *CronClient) CreateJob(req CreateJobRequest) (int64, error) {
	resp, err := c.doRequest("POST", "/api/cron-jobs", req)
	if err != nil {
		return 0, err
	}

	var result struct {
		ID int64 `json:"id"`
	}
	if err := json.Unmarshal(resp, &result); err != nil {
		return 0, fmt.Errorf("failed to parse response: %w", err)
	}
	return result.ID, nil
}

// ToggleJob enables or disables a job
func (c *CronClient) ToggleJob(id int64, enabled bool) error {
	_, err := c.doRequest("POST", fmt.Sprintf("/api/cron-jobs/%d/toggle", id), map[string]bool{"enabled": enabled})
	return err
}

// RunJob triggers a job immediately
func (c *CronClient) RunJob(id int64) error {
	_, err := c.doRequest("POST", fmt.Sprintf("/api/cron-jobs/%d/run", id), nil)
	return err
}

// GetJobLogs returns execution logs for a job
func (c *CronClient) GetJobLogs(id int64, limit int) ([]CronRunLog, error) {
	path := fmt.Sprintf("/api/cron-jobs/%d/logs?limit=%d", id, limit)
	resp, err := c.doRequest("GET", path, nil)
	if err != nil {
		return nil, err
	}

	var logs []CronRunLog
	if err := json.Unmarshal(resp, &logs); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}
	return logs, nil
}

// UpdateJobRequest is the request body for updating a job
type UpdateJobRequest struct {
	Name     *string `json:"name,omitempty"`
	Schedule *string `json:"schedule,omitempty"`
	TZ       *string `json:"tz,omitempty"`
	Message  *string `json:"message,omitempty"`
	Delivery *string `json:"delivery,omitempty"`
}

// UpdateJob updates a cron job
func (c *CronClient) UpdateJob(id int64, req UpdateJobRequest) error {
	_, err := c.doRequest("PUT", fmt.Sprintf("/api/cron-jobs/%d", id), req)
	return err
}

// DeleteJob deletes a cron job
func (c *CronClient) DeleteJob(id int64) error {
	_, err := c.doRequest("DELETE", fmt.Sprintf("/api/cron-jobs/%d", id), nil)
	return err
}
