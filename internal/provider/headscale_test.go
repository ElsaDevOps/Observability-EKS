package provider_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ElsaDevOps/Observability-EKS/internal/provider"
)

func TestCheckAPI_Success(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})
	server := httptest.NewServer(handler)
	defer server.Close()

	h := provider.NewHeadscale(server.URL, "api-key")
	ctx := context.Background()
	healthy, latency, err := h.CheckAPI(ctx)

	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
	if !healthy {
		t.Errorf("expected healthy=true, got %v", healthy)
	}
	if latency <= 0 {
		t.Errorf("expected latency > 0, got %v", latency)
	}
}

func TestCheckAPI_500(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	})
	server := httptest.NewServer(handler)
	defer server.Close()

	h := provider.NewHeadscale(server.URL, "api-key")
	ctx := context.Background()
	healthy, latency, err := h.CheckAPI(ctx)

	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
	if healthy {
		t.Errorf("expected healthy=false, got %v", healthy)
	}
	if latency <= 0 {
		t.Errorf("expected latency > 0, got %v", latency)
	}
}

func TestCheckAPI_Unreachable(t *testing.T) {
	h := provider.NewHeadscale("http://localhost:12345", "api-key")
	ctx := context.Background()
	healthy, latency, err := h.CheckAPI(ctx)

	if err == nil {
		t.Errorf("expected error, got %v", err)
	}
	if healthy {
		t.Errorf("expected healthy=false, got %v", healthy)
	}
	if latency != 0 {
		t.Errorf("expected latency = 0, got %v", latency)
	}
}
