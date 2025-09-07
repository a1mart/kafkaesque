package connectors

type Connector interface {
	Init(config map[string]string) error
	Start() error
	Stop() error
}
