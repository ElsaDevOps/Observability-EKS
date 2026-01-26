package provider

import (
	"context"
	"time"
)

type Node struct {
    Name     string `json:"name"`
    Online   bool `json:"online"`
    LastSeen time.Time `json:"last_seen"`
}

type listNodesResponse struct {
    Nodes []Node `json:"nodes"`
}


type Provider interface {
	Name() string
	CheckAPI(ctx context.Context) (healthy bool, latency time.Duration, err error)
	ListNodes(ctx context.Context) ([]Node, error)
}
