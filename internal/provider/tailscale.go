package provider

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

type tailscaleResponse struct {
	Devices []tailscaleDevice `json:"devices"`
}

type tailscaleDevice struct {
	Name               string    `json:"name"`
	ConnectedToControl bool      `json:"connectedToControl"`
	LastSeen           time.Time `json:"lastSeen"`
}

type Tailscale struct {
	APIKey    string
	TailnetID string
	URL       string
}

func NewTailscale(cfg ProviderConfig) Provider {
	url := "https://api.tailscale.com"
	if cfg.URL != "" {
		url = cfg.URL
	}

	return &Tailscale{TailnetID: cfg.TailnetID, APIKey: cfg.APIKey, URL: url}
}

func (t *Tailscale) Name() string {
	return "tailscale"
}

func (t *Tailscale) CheckAPI(ctx context.Context) (bool, time.Duration, error) {

	req, err := http.NewRequestWithContext(ctx, "GET", t.URL, nil)
	if err != nil {
		return false, 0, err
	}
	req.Header.Set("Authorization", "Bearer "+t.APIKey)

	start := time.Now()

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {

		return false, 0, err
	}
	defer resp.Body.Close()

	elapsed := time.Since(start)

	return resp.StatusCode == 200, elapsed, nil

}

func (t *Tailscale) ListNodes(ctx context.Context) ([]Node, error) {

	req, err := http.NewRequestWithContext(ctx, "GET", t.URL, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+t.APIKey)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {

		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("Expected 200, got: %s ", resp.Status)
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var response tailscaleResponse
	err = json.Unmarshal(data, &response)
	if err != nil {
		return nil, err
	}

	var nodes []Node
	for _, device := range response.Devices {
		nodes = append(nodes, Node{
			Name:     device.Name,
			Online:   device.ConnectedToControl,
			LastSeen: device.LastSeen,
		})
	}
	return nodes, nil

}

func init() {
	Register("tailscale", NewTailscale)
}
