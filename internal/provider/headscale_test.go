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

	h := provider.NewHeadscale(provider.ProviderConfig{
		URL:    server.URL,
		APIKey: "api-key",
	})
	ctx := context.Background()
	healthy, latency, err := h.CheckAPI(ctx)

	if err != nil {
		t.Errorf("got error=%v, want nil", err)
	}
	if !healthy {
		t.Errorf("got healthy=%v, want true", healthy)
	}
	if latency <= 0 {
		t.Errorf("got latency < %v, want latency > 0", latency)
	}
}

func TestCheckAPI_500(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	})
	server := httptest.NewServer(handler)
	defer server.Close()

	h := provider.NewHeadscale(provider.ProviderConfig{
		URL:    server.URL,
		APIKey: "api-key",
	})
	ctx := context.Background()
	healthy, latency, err := h.CheckAPI(ctx)

	if err != nil {
		t.Errorf("got error=%v, want nil", err)
	}
	if healthy {
		t.Errorf("got healthy=%v, want false", healthy)
	}
	if latency <= 0 {
		t.Errorf("got latency < %v, want latency > 0", latency)
	}
}

func TestCheckAPI_Unreachable(t *testing.T) {
	h := provider.NewHeadscale(provider.ProviderConfig{
		URL:    "http://localhost:12345",
		APIKey: "api-key",
	})
	ctx := context.Background()
	healthy, latency, err := h.CheckAPI(ctx)

	if err == nil {
		t.Errorf("got error=%v, want non-nil", err)
	}
	if healthy {
		t.Errorf("got healthy=%v, want false", healthy)
	}
	if latency != 0 {
		t.Errorf("got latency=%v, want 0", latency)
	}
}

func TestListNodes_Success(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{
  "nodes": [
    {"name": "node1", "online": true, "last_seen": "2024-01-01T00:00:00Z"},
    {"name": "node2", "online": false, "last_seen": "2024-01-01T00:00:00Z"},
	{"name": "node3", "online": false, "last_seen": "2024-01-01T00:00:00Z"}
]
}`))
	})
	server := httptest.NewServer(handler)
	defer server.Close()
	h := provider.NewHeadscale(provider.ProviderConfig{
		URL:    server.URL,
		APIKey: "api-key",
	})
	ctx := context.Background()
	nodes, err := h.ListNodes(ctx)

	if err != nil {
		t.Errorf("got error=%v, want nil", err)
	}
	if len(nodes) != 3 {
		t.Errorf("got node length=%v, want 3", len(nodes))
	}
}

func TestListNodes_500(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	})
	server := httptest.NewServer(handler)
	defer server.Close()

	h := provider.NewHeadscale(provider.ProviderConfig{
		URL:    server.URL,
		APIKey: "api-key",
	})
	ctx := context.Background()
	nodes, err := h.ListNodes(ctx)

	if err == nil {
		t.Errorf("got error=%v, want non-nil", nil)
	}
	if nodes != nil {
		t.Errorf("got nodes=%v, want 0", nodes)
	}
}

func TestListNodes_BadJSON(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{garbage`))

	})
	server := httptest.NewServer(handler)
	defer server.Close()

	h := provider.NewHeadscale(provider.ProviderConfig{
		URL:    server.URL,
		APIKey: "api-key",
	})
	ctx := context.Background()
	nodes, err := h.ListNodes(ctx)

	if err == nil {
		t.Errorf("got error=%v, want non-nil", nil)
	}
	if nodes != nil {
		t.Errorf("got nodes=%v, want 0", nodes)
	}
}
