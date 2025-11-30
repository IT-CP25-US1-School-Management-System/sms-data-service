package client

import (
	"context"

	"github.com/jmoiron/sqlx"
	"github.com/opentracing/opentracing-go"
)

// Client represents a database client that wraps sqlx.DB
type Client struct {
	db            *sqlx.DB
	connectionURI string
	driverName    string
	dbType        string // The logical database type (postgres, mysql, mssql, etc.)
	tracer        opentracing.Tracer
}

// NewClientFromFields creates a client from individual fields (used by adapters)
func NewClientFromFields(db *sqlx.DB, connectionURI, driverName, dbType string, tracer opentracing.Tracer) *Client {
	return &Client{
		db:            db,
		connectionURI: connectionURI,
		driverName:    driverName,
		dbType:        dbType,
		tracer:        tracer,
	}
}

// ClientConfig holds configuration for creating a new client
type ClientConfig struct {
	ConnectionString string
	DBType           string
	Tracer           opentracing.Tracer
}

// NewClient creates a new database client using the registered adapter
// This is imported from the parent database package to avoid circular dependency
var NewClient func(config ClientConfig) (*Client, error)

// NewClientWithContext creates a new database client with context
// This is imported from the parent database package to avoid circular dependency
var NewClientWithContext func(ctx context.Context, config ClientConfig) (*Client, error)

// GetClient returns the underlying sqlx.DB
func (c *Client) GetClient() *sqlx.DB {
	return c.db
}

// GetConnectionURI returns the connection URI
func (c *Client) GetConnectionURI() string {
	return c.connectionURI
}

// GetDriverName returns the driver name
func (c *Client) GetDriverName() string {
	return c.driverName
}

// GetDBType returns the logical database type
func (c *Client) GetDBType() string {
	return c.dbType
}

// SetDB sets the underlying sqlx.DB
func (c *Client) SetDB(db *sqlx.DB) {
	c.db = db
}

// IsConnect checks if the connection is alive
func (c *Client) IsConnect() bool {
	if err := c.db.Ping(); err == nil {
		return true
	}
	return false
}

// Close closes the database connection
func (c *Client) Close() error {
	return c.db.Close()
}
