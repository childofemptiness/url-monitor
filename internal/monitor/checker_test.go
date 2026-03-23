package monitor

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestCheckerChecker_OK(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer srv.Close()

	checker := Checker{}

	got := checker.Check(context.Background(), Monitor{
		URL: srv.URL,
		IntervalSeconds: 10,
	})

	if got.Status != MonitorCheckStatusUp {
		t.Fatalf("expected status %q, got %q", MonitorCheckStatusUp, got.Status)
	}

	if got.HTTPStatusCode != http.StatusOK {
		t.Fatalf("expected http status %q, got %q", http.StatusOK, got.HTTPStatusCode)
	}
}

func TestCheckerChecker_Down(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer srv.Close()

	checker := Checker{}

	got := checker.Check(context.Background(), Monitor{
		URL: srv.URL,
		IntervalSeconds: 10,
	})

	if got.Status != MonitorCheckStatusDown {
		t.Fatalf("expected status %q, got %q", MonitorCheckStatusDown, got.Status)
	}

	if got.HTTPStatusCode != http.StatusInternalServerError {
		t.Fatalf("expected http status %q, got %q", http.StatusInternalServerError, got.HTTPStatusCode)
	}
}

func TestCheckerChecker_Error(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	srv.Close()

	checker := Checker{}

	got := checker.Check(context.Background(), Monitor{
		URL: srv.URL,
		IntervalSeconds: 10,
	})

	if got.Status != MonitorCheckStatusError {
		t.Fatalf("expected status %q, got %q", MonitorCheckStatusError, got.Status)
	}

	if got.ErrorMessage == "" {
		t.Fatal("expected non-empty error message")
	}
}
