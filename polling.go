package bizverify

import (
	"context"
	"fmt"
	"time"
)

const (
	defaultPollInterval = 2 * time.Second
	maxPollInterval     = 10 * time.Second
	defaultPollTimeout  = 60 * time.Second
)

func pollUntilComplete(ctx context.Context, c *httpClient, resp *VerifyResponse, opts *PollOptions) (*JobStatusResponse, error) {
	pollInterval := defaultPollInterval
	timeout := defaultPollTimeout
	var onStatusChange func(JobStatusResponse)

	if opts != nil {
		if opts.PollInterval > 0 {
			pollInterval = opts.PollInterval
		}
		if opts.Timeout > 0 {
			timeout = opts.Timeout
		}
		onStatusChange = opts.OnStatusChange
	}

	if resp.Status == "completed" && resp.JobID != nil {
		var status JobStatusResponse
		err := c.request(ctx, requestOptions{
			method: "GET",
			path:   "/v1/verify/status/" + *resp.JobID,
			auth:   authAPIKey,
		}, &status)
		if err != nil {
			return nil, err
		}
		return &status, nil
	}

	if resp.JobID == nil {
		return nil, fmt.Errorf("bizverify: no job_id in verify response — cannot poll")
	}
	jobID := *resp.JobID

	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	interval := pollInterval
	var lastStatus string

	for {
		var status JobStatusResponse
		err := c.request(ctx, requestOptions{
			method: "GET",
			path:   "/v1/verify/status/" + jobID,
			auth:   authAPIKey,
		}, &status)
		if err != nil {
			if ctx.Err() != nil {
				return nil, &TimeoutError{APIError: &APIError{
					Message: fmt.Sprintf("Polling timed out after %s for job %s", timeout, jobID),
					Code:    "TIMEOUT",
				}}
			}
			return nil, err
		}

		if onStatusChange != nil && status.Status != lastStatus {
			lastStatus = status.Status
			onStatusChange(status)
		}

		if status.Status == "completed" {
			return &status, nil
		}

		if status.Status == "failed" {
			errMsg := "Verification job failed"
			if status.Error != nil {
				errMsg = *status.Error
			}
			return nil, &JobFailedError{
				APIError: &APIError{Message: errMsg, Code: "JOB_FAILED"},
				JobID: jobID,
			}
		}

		select {
		case <-time.After(interval):
		case <-ctx.Done():
			return nil, &TimeoutError{APIError: &APIError{
				Message: fmt.Sprintf("Polling timed out after %s for job %s", timeout, jobID),
				Code:    "TIMEOUT",
			}}
		}

		newInterval := time.Duration(float64(interval) * 1.5)
		if newInterval > maxPollInterval {
			newInterval = maxPollInterval
		}
		interval = newInterval
	}
}
