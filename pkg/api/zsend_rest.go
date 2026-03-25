package api

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/zeabur/cli/pkg/constant"
	"github.com/zeabur/cli/pkg/model"
)

const zsendRESTBase = constant.ZeaburServerURL + "/api/v1/zsend"

func zsendDo(ctx context.Context, apiKey, method, path string, body any) ([]byte, int, error) {
	var bodyReader io.Reader
	if body != nil {
		b, err := json.Marshal(body)
		if err != nil {
			return nil, 0, fmt.Errorf("marshal request: %w", err)
		}
		bodyReader = bytes.NewReader(b)
	}

	req, err := http.NewRequestWithContext(ctx, method, zsendRESTBase+path, bodyReader)
	if err != nil {
		return nil, 0, fmt.Errorf("create request: %w", err)
	}
	req.Header.Set("Authorization", "Bearer "+apiKey)
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	c := &http.Client{Timeout: 30 * time.Second}
	resp, err := c.Do(req)
	if err != nil {
		return nil, 0, fmt.Errorf("do request: %w", err)
	}
	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, resp.StatusCode, fmt.Errorf("read response: %w", err)
	}

	if resp.StatusCode >= 400 {
		var errResp struct {
			Error string `json:"error"`
		}
		if jsonErr := json.Unmarshal(data, &errResp); jsonErr == nil && errResp.Error != "" {
			return nil, resp.StatusCode, fmt.Errorf("API error (%d): %s", resp.StatusCode, errResp.Error)
		}
		return nil, resp.StatusCode, fmt.Errorf("API error (%d)", resp.StatusCode)
	}

	return data, resp.StatusCode, nil
}

func (c *client) SendZSendEmail(ctx context.Context, apiKey string, req model.ZSendSendEmailRequest) (*model.ZSendSendEmailReply, error) {
	data, _, err := zsendDo(ctx, apiKey, http.MethodPost, "/emails", req)
	if err != nil {
		return nil, err
	}
	var reply model.ZSendSendEmailReply
	if err := json.Unmarshal(data, &reply); err != nil {
		return nil, fmt.Errorf("decode response: %w", err)
	}
	return &reply, nil
}

func (c *client) ScheduleZSendEmail(ctx context.Context, apiKey string, req model.ZSendScheduleEmailRequest) (*model.ZSendScheduleEmailReply, error) {
	data, _, err := zsendDo(ctx, apiKey, http.MethodPost, "/emails/schedule", req)
	if err != nil {
		return nil, err
	}
	var reply model.ZSendScheduleEmailReply
	if err := json.Unmarshal(data, &reply); err != nil {
		return nil, fmt.Errorf("decode response: %w", err)
	}
	return &reply, nil
}

func (c *client) SendZSendBatchEmail(ctx context.Context, apiKey string, req model.ZSendBatchEmailRequest) (*model.ZSendBatchEmailReply, error) {
	data, _, err := zsendDo(ctx, apiKey, http.MethodPost, "/emails/batch", req)
	if err != nil {
		return nil, err
	}
	var reply model.ZSendBatchEmailReply
	if err := json.Unmarshal(data, &reply); err != nil {
		return nil, fmt.Errorf("decode response: %w", err)
	}
	return &reply, nil
}

func (c *client) ListZSendScheduledEmails(ctx context.Context, apiKey string, page, pageSize *int, status *string) (*model.ZSendListScheduledEmailsReply, error) {
	q := url.Values{}
	if page != nil {
		q.Set("page", strconv.Itoa(*page))
	}
	if pageSize != nil {
		q.Set("page_size", strconv.Itoa(*pageSize))
	}
	if status != nil && *status != "" {
		q.Set("status", *status)
	}
	path := "/emails/scheduled"
	if len(q) > 0 {
		path += "?" + q.Encode()
	}

	data, _, err := zsendDo(ctx, apiKey, http.MethodGet, path, nil)
	if err != nil {
		return nil, err
	}
	var reply model.ZSendListScheduledEmailsReply
	if err := json.Unmarshal(data, &reply); err != nil {
		return nil, fmt.Errorf("decode response: %w", err)
	}
	return &reply, nil
}

func (c *client) GetZSendScheduledEmail(ctx context.Context, apiKey string, id string) (*model.ZSendScheduledEmail, error) {
	data, _, err := zsendDo(ctx, apiKey, http.MethodGet, "/emails/scheduled/"+id, nil)
	if err != nil {
		return nil, err
	}
	var reply model.ZSendScheduledEmail
	if err := json.Unmarshal(data, &reply); err != nil {
		return nil, fmt.Errorf("decode response: %w", err)
	}
	return &reply, nil
}

func (c *client) CancelZSendScheduledEmail(ctx context.Context, apiKey string, id string) error {
	_, _, err := zsendDo(ctx, apiKey, http.MethodDelete, "/emails/scheduled/"+id, nil)
	return err
}

func (c *client) ListZSendBatchEmailJobs(ctx context.Context, apiKey string, page, pageSize *int, status *string) (*model.ZSendListBatchJobsReply, error) {
	q := url.Values{}
	if page != nil {
		q.Set("page", strconv.Itoa(*page))
	}
	if pageSize != nil {
		q.Set("page_size", strconv.Itoa(*pageSize))
	}
	if status != nil && *status != "" {
		q.Set("status", *status)
	}
	path := "/emails/batch"
	if len(q) > 0 {
		path += "?" + q.Encode()
	}

	data, _, err := zsendDo(ctx, apiKey, http.MethodGet, path, nil)
	if err != nil {
		return nil, err
	}
	var reply model.ZSendListBatchJobsReply
	if err := json.Unmarshal(data, &reply); err != nil {
		return nil, fmt.Errorf("decode response: %w", err)
	}
	return &reply, nil
}

func (c *client) GetZSendBatchEmailJob(ctx context.Context, apiKey string, id string) (*model.ZSendBatchJob, error) {
	data, _, err := zsendDo(ctx, apiKey, http.MethodGet, "/emails/batch/"+id, nil)
	if err != nil {
		return nil, err
	}
	var reply model.ZSendBatchJob
	if err := json.Unmarshal(data, &reply); err != nil {
		return nil, fmt.Errorf("decode response: %w", err)
	}
	return &reply, nil
}
