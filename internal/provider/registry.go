package provider

type FactoryFunc func(ProviderConfig) Provider

var registry = map[string]FactoryFunc{}

func Register(key string, value FactoryFunc) {
	registry[key] = value
}

func GetRegistry(name string) (FactoryFunc, bool) {
	factory, ok := registry[name]
	return factory, ok
}
