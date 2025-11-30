package database

import (
	"context"
	"fmt"
	"sync"

	"github.com/IT-CP25-US1-School-Management-System/sms-data-service/models/entity"
	"github.com/IT-CP25-US1-School-Management-System/sms-data-service/service/database/v1/client"
	"github.com/Masterminds/squirrel"
)

// DatabaseAdapter defines the complete interface for database support
// Each database type implements this single interface for both connection and SQL building
type DatabaseAdapter interface {
	// GetName returns the database type name (e.g., "postgres", "mysql", "mssql")
	GetName() string

	// Connect creates a new database connection
	Connect(ctx context.Context, config client.ClientConfig) (*client.Client, error)

	// GetDriverName returns the SQL driver name used by database/sql
	GetDriverName() string

	// BuildConnectionString creates a connection string from source configuration
	BuildConnectionString(source *entity.Sources, decryptedPassword string) (string, error)

	// GetPlaceholderFormat returns the SQL placeholder format for this database
	// Examples: PostgreSQL uses $1, MySQL uses ?, SQL Server uses @p1
	GetPlaceholderFormat() squirrel.PlaceholderFormat

	// GetDriverType returns the driver type for database connection
	GetDriverType() string

	// BuildJSONObject creates a JSON object from key-value pairs
	// keyValuePairs should be alternating ['key1', value1, 'key2', value2, ...]
	// Returns database-specific JSON object construction syntax
	BuildJSONObject(keyValuePairs []string) string

	// BuildJSONArrayAgg creates a JSON array aggregation in a subquery
	// Returns database-specific JSON array aggregation syntax
	BuildJSONArrayAgg(jsonObjectSQL, table, alias, whereClause, resultAlias string) string

	// GetInformationTablesQuery returns the query to fetch table information
	// Returns database-specific query for information_schema tables
	GetInformationTablesQuery() string

	// GetInformationColumnsQuery returns the query to fetch column information
	// Returns database-specific query for information_schema columns
	GetInformationColumnsQuery() string

	// GetInformationSchemasQuery returns the query to fetch schema information
	// Returns database-specific query for information_schema schemata
	GetInformationSchemasQuery() string

	// GetInformationTableRelationsQuery returns the query to fetch table relationships (foreign keys)
	// Returns database-specific query for information_schema foreign key relationships
	GetInformationTableRelationsQuery() string
}

// AdapterRegistry manages all registered database adapters
type AdapterRegistry struct {
	adapters map[string]DatabaseAdapter
	mu       sync.RWMutex
}

// NewAdapterRegistry creates a new adapter registry
func NewAdapterRegistry() *AdapterRegistry {
	return &AdapterRegistry{
		adapters: make(map[string]DatabaseAdapter),
	}
}

// Register registers a new database adapter
func (r *AdapterRegistry) Register(adapter DatabaseAdapter) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	name := adapter.GetName()
	if _, exists := r.adapters[name]; exists {
		return fmt.Errorf("adapter '%s' is already registered", name)
	}
	r.adapters[name] = adapter
	return nil
}

// Get retrieves an adapter by database type
func (r *AdapterRegistry) Get(dbType string) (DatabaseAdapter, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	adapter, exists := r.adapters[dbType]
	if !exists {
		return nil, fmt.Errorf("unsupported database type: %s", dbType)
	}
	return adapter, nil
}

// List returns all registered adapter names
func (r *AdapterRegistry) List() []string {
	r.mu.RLock()
	defer r.mu.RUnlock()

	names := make([]string, 0, len(r.adapters))
	for name := range r.adapters {
		names = append(names, name)
	}
	return names
}

// Global adapter registry instance
var globalAdapterRegistry = NewAdapterRegistry()

// RegisterAdapter registers an adapter to the global registry
func RegisterAdapter(adapter DatabaseAdapter) error {
	return globalAdapterRegistry.Register(adapter)
}

// GetAdapter retrieves an adapter from the global registry
func GetAdapter(dbType string) (DatabaseAdapter, error) {
	return globalAdapterRegistry.Get(dbType)
}

// ListAdapters returns all registered adapter names
func ListAdapters() []string {
	return globalAdapterRegistry.List()
}

// NewClient creates a new database client using the registered adapter
func NewClient(config client.ClientConfig) (*client.Client, error) {
	adapter, err := GetAdapter(config.DBType)
	if err != nil {
		return nil, fmt.Errorf("failed to get adapter for %s: %w", config.DBType, err)
	}

	return adapter.Connect(context.Background(), config)
}

// NewClientWithContext creates a new database client with context
func NewClientWithContext(ctx context.Context, config client.ClientConfig) (*client.Client, error) {
	adapter, err := GetAdapter(config.DBType)
	if err != nil {
		return nil, fmt.Errorf("failed to get adapter for %s: %w", config.DBType, err)
	}

	return adapter.Connect(ctx, config)
}

func init() {
	// Inject the NewClient functions into client package to avoid circular dependency
	client.NewClient = NewClient
	client.NewClientWithContext = NewClientWithContext
}
