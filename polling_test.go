package bizverify

import (
	"context"
	"errors"
	"net/http"
	"sync/atomic"
	"testing"
	"time"
)

func TestPollAlreadyCompleted(t *testing.T) {
	client, _ := setupTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, 200, fixtureJobCompletedResp)
	})
	jobID := "job_456"
	resp := &VerifyResponse{Status: "completed", JobID: &jobID, Cached: true, CreditsCharged: 1}
	result, err := pollUntilComplete(context.Background(), client.client, resp, &PollOptions{
		PollInterval: time.Millisecond, Timeout: time.Second,
	})
	if err != nil {
		t.Fatal(err)
	}
	if result.Status != "completed" {
		t.Errorf("expected completed, got %s", result.Status)
	}
}

func TestPollPendingToCompleted(t *testing.T) {
	var count int64
	client, _ := setupTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		n := atomic.AddInt64(&count, 1)
		if n < 2 {
			writeJSON(w, 200, fixtureJobPendingResp)
			return
		}
		writeJSON(w, 200, fixtureJobCompletedResp)
	})
	jobID := "job_456"
	resp := &VerifyResponse{Status: "pending", JobID: &jobID}
	result, err := pollUntilComplete(context.Background(), client.client, resp, &PollOptions{
		PollInterval: time.Millisecond, Timeout: 5 * time.Second,
	})
	if err != nil {
		t.Fatal(err)
	}
	if result.Status != "completed" {
		t.Errorf("expected completed, got %s", result.Status)
	}
}

func TestPollPendingToFailed(t *testing.T) {
	var count int64
	client, _ := setupTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		n := atomic.AddInt64(&count, 1)
		if n < 2 {
			writeJSON(w, 200, fixtureJobPendingResp)
			return
		}
		writeJSON(w, 200, fixtureJobFailedResp)
	})
	jobID := "job_456"
	resp := &VerifyResponse{Status: "pending", JobID: &jobID}
	_, err := pollUntilComplete(context.Background(), client.client, resp, &PollOptions{
		PollInterval: time.Millisecond, Timeout: 5 * time.Second,
	})
	var jfe *JobFailedError
	if !errors.As(err, &jfe) {
		t.Fatalf("expected JobFailedError, got %T: %v", err, err)
	}
	if jfe.JobID != "job_456" {
		t.Errorf("wrong job ID: %s", jfe.JobID)
	}
}

func TestPollTimeout(t *testing.T) {
	client, _ := setupTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, 200, fixtureJobPendingResp)
	})
	jobID := "job_456"
	resp := &VerifyResponse{Status: "pending", JobID: &jobID}
	_, err := pollUntilComplete(context.Background(), client.client, resp, &PollOptions{
		PollInterval: time.Millisecond, Timeout: 50 * time.Millisecond,
	})
	var te *TimeoutError
	if !errors.As(err, &te) {
		t.Fatalf("expected TimeoutError, got %T: %v", err, err)
	}
}

func TestPollNoJobID(t *testing.T) {
	client := New()
	resp := &VerifyResponse{Status: "pending"}
	_, err := pollUntilComplete(context.Background(), client.client, resp, nil)
	if err == nil {
		t.Fatal("expected error for missing job_id")
	}
}
