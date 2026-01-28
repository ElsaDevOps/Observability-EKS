package provider

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

type Headscale struct {
	url    string
	apikey string
}

func (h *Headscale) ListNodes(ctx context.Context) ([]Node, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", h.url+"/api/v1/node", nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+h.apikey)

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

	var response listNodesResponse
	err = json.Unmarshal(data, &response)
	if err != nil {
		return nil, err
	}
	return response.Nodes, nil
}

func (h *Headscale) Name() string {
	return "headscale"
}

func (h *Headscale) CheckAPI(ctx context.Context) (bool, time.Duration, error) {

	req, err := http.NewRequestWithContext(ctx, "GET", h.url+"/api/v1/apikey", nil)
	if err != nil {
		return false, 0, err
	}
	req.Header.Set("Authorization", "Bearer "+h.apikey)

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

func NewHeadscale(url, apikey string) *Headscale {
	return &Headscale{url: url, apikey: apikey}
}
