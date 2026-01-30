package provider

import (
	"context"
	"time"
)

type Node struct {
	Name string
	//Online in the Node struct represents whether a node is connected/online in the network,
	// not whether the API itself is healthy
	Online   bool
	LastSeen time.Time
}

type listNodesResponse struct {
	Nodes []Node `json:"nodes"`
}

type Provider interface {
	Name() string
	CheckAPI(ctx context.Context) (healthy bool, latency time.Duration, err error)
	ListNodes(ctx context.Context) ([]Node, error)
}

type ProviderConfig struct {
	Name      string        `yaml:"name"`
	URL       string        `yaml:"url"`
	APIKey    string        `yaml:"api_key"`
	TailnetID string        `yaml:"tailnet_id"`
	Interval  time.Duration // optional, overrides DefaultInterval
	Features  map[string]string
}
