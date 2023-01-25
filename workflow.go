package onfido

import (
	"bytes"
	"context"
	"encoding/json"
	"time"
)

type WorkflowRun struct {
	WorkflowID  string `json:"workflow_id"`
	ApplicantID string `json:"applicant_id"`
}

type Workflow struct {
	ApplicantID  string      `json:"applicant_id"`
	CreatedAt    time.Time   `json:"created_at"`
	DashboardURL string      `json:"dashboard_url"`
	Error        interface{} `json:"error"`
	ID           string      `json:"id"`
	Link         struct {
		CompletedRedirectURL interface{} `json:"completed_redirect_url"`
		ExpiredRedirectURL   interface{} `json:"expired_redirect_url"`
		ExpiresAt            interface{} `json:"expires_at"`
		Language             interface{} `json:"language"`
		URL                  string      `json:"url"`
	} `json:"link"`
	Output            interface{}   `json:"output"`
	Reasons           []interface{} `json:"reasons"`
	Status            string        `json:"status"`
	UpdatedAt         time.Time     `json:"updated_at"`
	WorkflowID        string        `json:"workflow_id"`
	WorkflowVersionID int           `json:"workflow_version_id"`
}

func (c *Client) CreateWorkflowRun(ctx context.Context, a WorkflowRun) (*Workflow, error) {
	jsonStr, err := json.Marshal(a)
	if err != nil {
		return nil, err
	}

	req, err := c.newRequest("POST", "/workflow_runs", bytes.NewBuffer(jsonStr))
	if err != nil {
		return nil, err
	}

	var resp Workflow
	_, err = c.do(ctx, req, &resp)
	return &resp, err
}

func (c *Client) GetWorkflowRun(ctx context.Context, id string) (*Workflow, error) {
	req, err := c.newRequest("GET", "/workflow_runs/"+id, nil)
	if err != nil {
		return nil, err
	}

	var resp Workflow
	_, err = c.do(ctx, req, &resp)
	return &resp, err
}
