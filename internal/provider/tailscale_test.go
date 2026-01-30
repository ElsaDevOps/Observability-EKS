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

	ts := provider.NewTailscale(provider.ProviderConfig{
		URL:       server.URL,
		TailnetID: "test-tailnet",
		APIKey:    "test-key",
	})
	ctx := context.Background()
	healthy, latency, err := ts.CheckAPI(ctx)

	if err != nil {
		t.Errorf("got err=%v, want nil", err)
	}
	if !healthy {
		t.Errorf("got healthy=%v, want true", healthy)
	}
	if latency <= 0 {
		t.Errorf("got latency=%v, want >0", latency)
	}
}

func TestCheckAPI_500_TS(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	})
	server := httptest.NewServer(handler)
	defer server.Close()
	ts := provider.NewTailscale(provider.ProviderConfig{
		URL:       server.URL,
		TailnetID: "test-tailnet",
		APIKey:    "test-key",
	})
	ctx := context.Background()
	healthy, latency, err := ts.CheckAPI(ctx)

	if err != nil {
		t.Errorf("got err=%v, want nil", err)
	}
	if healthy {
		t.Errorf("got healthy=%v, want false", healthy)
	}
	if latency <= 0 {
		t.Errorf("got latency=%v, want >0", latency)
	}
}

func TestCheckAPI_Unreachable_TS(t *testing.T) {
	ts := provider.NewTailscale(provider.ProviderConfig{
		URL:       "http://localhost:12345",
		TailnetID: "test-tailnet",
		APIKey:    "test-key",
	})
	ctx := context.Background()
	healthy, latency, err := ts.CheckAPI(ctx)

	if err == nil {
		t.Errorf("got err=nil, want non-nil")
	}
	if healthy {
		t.Errorf("got healthy=%v, want false", healthy)
	}
	if latency != 0 {
		t.Errorf("got latency=%v, want 0", latency)
	}
}

func TestListNodes_Success_TS(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"devices": [{"name": "test-node", "connectedToControl": true, "lastSeen": "2024-01-01T00:00:00Z"}, {"name": "test-node-2", "connectedToControl": false, "lastSeen": "2024-01-02T00:00:00Z"}, {"name": "test-node-3", "connectedToControl": true, "lastSeen": "2024-01-03T00:00:00Z"}]}`))
	})
	server := httptest.NewServer(handler)
	defer server.Close()

	ts := provider.NewTailscale(provider.ProviderConfig{
		URL:       server.URL,
		TailnetID: "test-tailnet",
		APIKey:    "test-key",
	})
	ctx := context.Background()
	nodes, err := ts.ListNodes(ctx)

	if err != nil {
		t.Errorf("got err=%v, want nil", err)
	}
	if len(nodes) != 3 {
		t.Errorf("got len(nodes)=%v, want 3", len(nodes))
	}
	if nodes[0].Name != "test-node" {
		t.Errorf("got name=%s, want test-node", nodes[0].Name)
	}
	if !nodes[0].Online {
		t.Errorf("got Online=%v, want true", nodes[0].Online)
	}
}

func TestListNodes_500_TS(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	})
	server := httptest.NewServer(handler)
	defer server.Close()

	ts := provider.NewTailscale(provider.ProviderConfig{
		URL:       server.URL,
		TailnetID: "test-tailnet",
		APIKey:    "test-key",
	})
	ctx := context.Background()
	nodes, err := ts.ListNodes(ctx)

	if err == nil {
		t.Errorf("got err=nil, want non-nil")
	}
	if nodes != nil {
		t.Errorf("got nodes=%v, want nil", nodes)
	}
}

func TestListNodes_BadJSON_TS(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{garbage`))
	})
	server := httptest.NewServer(handler)
	defer server.Close()
	ts := provider.NewTailscale(provider.ProviderConfig{
		URL:       server.URL,
		TailnetID: "test-tailnet",
		APIKey:    "test-key",
	})
	ctx := context.Background()
	nodes, err := ts.ListNodes(ctx)

	if err == nil {
		t.Errorf("got err=nil, want non-nil")
	}
	if nodes != nil {
		t.Errorf("got nodes=%v, want nil", nodes)
	}
}
