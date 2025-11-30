package client

import (
	"context"
	"fmt"
	"sync"
)

// Connector defines the interface for database connectors
// Each database type implements this to handle its specific connection logic
type Connector interface {
	// GetName returns the database type name (e.g., "postgres", "mysql")
	GetName() string

	// Connect creates a new database connection
	Connect(ctx context.Context, config ClientConfig) (*Client, error)

	// GetDriverName returns the SQL driver name used by database/sql
	GetDriverName() string
}

// ConnectorRegistry manages all registered database connectors
type ConnectorRegistry struct {
	connectors map[string]Connector
	mu         sync.RWMutex
}

// NewConnectorRegistry creates a new connector registry
func NewConnectorRegistry() *ConnectorRegistry {
	return &ConnectorRegistry{
		connectors: make(map[string]Connector),
	}
}

// Register registers a new database connector
func (r *ConnectorRegistry) Register(connector Connector) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	name := connector.GetName()
	if _, exists := r.connectors[name]; exists {
		return fmt.Errorf("connector '%s' is already registered", name)
	}
	r.connectors[name] = connector
	return nil
}

// Get retrieves a connector by database type
func (r *ConnectorRegistry) Get(dbType string) (Connector, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	connector, exists := r.connectors[dbType]
	if !exists {
		return nil, fmt.Errorf("unsupported database type: %s", dbType)
	}
	return connector, nil
}

// List returns all registered connector names
func (r *ConnectorRegistry) List() []string {
	r.mu.RLock()
	defer r.mu.RUnlock()

	names := make([]string, 0, len(r.connectors))
	for name := range r.connectors {
		names = append(names, name)
	}
	return names
}

// Global connector registry instance
var globalConnectorRegistry = NewConnectorRegistry()

// RegisterConnector registers a connector to the global registry
func RegisterConnector(connector Connector) error {
	return globalConnectorRegistry.Register(connector)
}

// GetConnector retrieves a connector from the global registry
func GetConnector(dbType string) (Connector, error) {
	return globalConnectorRegistry.Get(dbType)
}

// ListConnectors returns all registered connector names
func ListConnectors() []string {
	return globalConnectorRegistry.List()
}
