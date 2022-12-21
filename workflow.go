package onfido

import (
	"bytes"
	"context"
	"encoding/json"
)

type WorkflowRun struct {
	WorkflowId  string `json:"workflow_id"`
	ApplicantId string `json:"applicant_id"`
}

type Workflow struct {
	Id string `json:"id"`
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
