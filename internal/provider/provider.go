package provider

type Node struct {
    ID      string
    Name    string
    Healthy bool
}

type Provider interface {
    Name() string
    CheckAPI(ctx context.Context) (healthy bool, latency time.Duration, err error)
    ListNodes(ctx context.Context) ([]Node, error)
}