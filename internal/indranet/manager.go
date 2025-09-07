// connectors/manager.go
package connectors

import (
	"errors"
	"sync"
)

type ConnectorManager struct {
	connectors map[string]Connector
	mu         sync.Mutex
}

func NewConnectorManager() *ConnectorManager {
	return &ConnectorManager{
		connectors: make(map[string]Connector),
	}
}

func (cm *ConnectorManager) Register(name string, connector Connector, config map[string]string) error {
	cm.mu.Lock()
	defer cm.mu.Unlock()

	if _, exists := cm.connectors[name]; exists {
		return errors.New("connector already registered")
	}

	if err := connector.Init(config); err != nil {
		return err
	}

	cm.connectors[name] = connector
	return nil
}

func (cm *ConnectorManager) Start(name string) error {
	cm.mu.Lock()
	defer cm.mu.Unlock()

	connector, exists := cm.connectors[name]
	if !exists {
		return errors.New("connector not found")
	}

	return connector.Start()
}

func (cm *ConnectorManager) Stop(name string) error {
	cm.mu.Lock()
	defer cm.mu.Unlock()

	connector, exists := cm.connectors[name]
	if !exists {
		return errors.New("connector not found")
	}

	err := connector.Stop()
	delete(cm.connectors, name)
	return err
}

func (cm *ConnectorManager) ListConnectors() []string {
	cm.mu.Lock()
	defer cm.mu.Unlock()

	var keys []string
	for k := range cm.connectors {
		keys = append(keys, k)
	}
	return keys
}
