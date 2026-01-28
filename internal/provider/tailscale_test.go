package provider_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ElsaDevOps/Observability-EKS/internal/provider"
)

func TestCheckAPI_Success_TS(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})
	server := httptest.NewServer(handler)
	defer server.Close()

	ts := &provider.Tailscale{
		URL:       server.URL,
		TailnetID: "test-tailnet",
		APIKey:    "test-key",
	}
	ctx := context.Background()
	healthy, latency, err := ts.CheckAPI(ctx)

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

func TestCheckAPI_500_TS(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	})
	server := httptest.NewServer(handler)
	defer server.Close()

	ts := &provider.Tailscale{
		URL:       server.URL,
		TailnetID: "test-tailnet",
		APIKey:    "test-key",
	}
	ctx := context.Background()
	healthy, latency, err := ts.CheckAPI(ctx)

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

func TestCheckAPI_Unreachable_TS(t *testing.T) {
	ts := &provider.Tailscale{
		URL:       "http://localhost:12345",
		TailnetID: "test-tailnet",
		APIKey:    "test-key",
	}
	ctx := context.Background()
	healthy, latency, err := ts.CheckAPI(ctx)

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

func TestListNodes_Success_TS(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"devices": [{"name": "test-node", "connectedToControl": true, "lastSeen": "2024-01-01T00:00:00Z"}, {"name": "test-node-2", "connectedToControl": false, "lastSeen": "2024-01-02T00:00:00Z"}, {"name": "test-node-3", "connectedToControl": true, "lastSeen": "2024-01-03T00:00:00Z"}]}`))

	})
	server := httptest.NewServer(handler)
	defer server.Close()

	ts := &provider.Tailscale{
		URL:       server.URL,
		TailnetID: "test-tailnet",
		APIKey:    "test-key",
	}
	ctx := context.Background()
	nodes, err := ts.ListNodes(ctx)

	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
	if len(nodes) != 3 {
		t.Errorf("expected node length of 3, got %v", len(nodes))
	}

	if nodes[0].Name != "test-node" {
		t.Errorf("expected name 'test-node', got %s", nodes[0].Name)
	}
	if !nodes[0].Online {
		t.Errorf("expected Online to be true")
	}

}

func TestListNodes_500_TS(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)

	})
	server := httptest.NewServer(handler)
	defer server.Close()

	ts := &provider.Tailscale{
		URL:       server.URL,
		TailnetID: "test-tailnet",
		APIKey:    "test-key",
	}

	ctx := context.Background()
	nodes, err := ts.ListNodes(ctx)

	if err == nil {
		t.Errorf("expected error, got %v", nil)
	}
	if nodes != nil {
		t.Errorf("expected node length of 0, got %v", nodes)
	}

}

func TestListNodes_BadJSON_TS(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{garbage`))

	})
	server := httptest.NewServer(handler)
	defer server.Close()
	ts := &provider.Tailscale{
		URL:       server.URL,
		TailnetID: "test-tailnet",
		APIKey:    "test-key",
	}
	ctx := context.Background()
	nodes, err := ts.ListNodes(ctx)

	if err == nil {
		t.Errorf("expected error, got %v", nil)
	}
	if nodes != nil {
		t.Errorf("expected node length of 0, got %v", nodes)
	}

}
