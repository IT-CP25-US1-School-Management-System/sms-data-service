package repository

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/IT-CP25-US1-School-Management-System/sms-data-service/models/entity"
	"github.com/IT-CP25-US1-School-Management-System/sms-data-service/service/data/v1"
	database "github.com/IT-CP25-US1-School-Management-System/sms-data-service/service/database/v1"
	"github.com/IT-CP25-US1-School-Management-System/sms-data-service/service/database/v1/client"
	"github.com/IT-CP25-US1-School-Management-System/sms-data-service/utils/crypto"
	"github.com/gofrs/uuid"
)

type dbConnectionManagerRepository struct {
	datasetRepo  data.PsqlDatasetRepository
	connections  map[uuid.UUID]*client.Client
	dbTypes      map[uuid.UUID]string // Store dbType for each connection
	mu           sync.RWMutex
	cryptoSecret string
}

func NewDBConnectionManagerRepository(datasetRepo data.PsqlDatasetRepository, cryptoSecret string) database.DBConnectionManagerRepository {
	return &dbConnectionManagerRepository{
		datasetRepo:  datasetRepo,
		cryptoSecret: cryptoSecret,
		connections:  make(map[uuid.UUID]*client.Client),
		dbTypes:      make(map[uuid.UUID]string),
	}
}

func (cm *dbConnectionManagerRepository) createConnectionDetails(source *entity.Sources) (string, error) {
	decryptedPass, err := crypto.Decrypt(source.Password, cm.cryptoSecret)
	if err != nil {
		return "", fmt.Errorf("failed to decrypt")
	}

	// Get adapter from registry
	adapter, err := database.GetAdapter(source.DBType)
	if err != nil {
		return "", err
	}

	// Build connection string using adapter
	connStr, err := adapter.BuildConnectionString(source, decryptedPass)
	if err != nil {
		return "", fmt.Errorf("failed to build connection string: %w", err)
	}

	return connStr, nil
}

func (cm *dbConnectionManagerRepository) GetConnection(ctx context.Context, sourceID uuid.UUID) (*client.Client, error) {
	cm.mu.RLock()
	if client, ok := cm.connections[sourceID]; ok {
		cm.mu.RUnlock()
		return client, nil
	}
	cm.mu.RUnlock()

	ctxFetch, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()
	source, err := cm.datasetRepo.FetchSourceByID(ctxFetch, &sourceID)
	if err != nil {
		return nil, fmt.Errorf("failed to get source: %w", err)
	}
	if source == nil {
		return nil, fmt.Errorf("source not found")
	}

	connStr, err := cm.createConnectionDetails(source)
	if err != nil {
		return nil, err
	}

	// Get adapter and use it to connect
	adapter, err := database.GetAdapter(source.DBType)
	if err != nil {
		return nil, fmt.Errorf("failed to get adapter: %w", err)
	}

	clientConn, err := adapter.Connect(ctx, client.ClientConfig{
		ConnectionString: connStr,
		DBType:           source.DBType,
		Tracer:           nil, // TODO: Add tracer if needed
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create connection: %w", err)
	}

	// ใส่ cache โดยตรวจซ้ำ
	cm.mu.Lock()
	if existing, ok := cm.connections[sourceID]; ok {
		cm.mu.Unlock()
		// มีคนสร้างทันก่อน ใช้อันเดิมและปิดอันใหม่
		clientConn.GetClient().Close()
		return existing, nil
	}
	cm.connections[sourceID] = clientConn
	cm.dbTypes[sourceID] = source.DBType // Store the dbType
	cm.mu.Unlock()
	return clientConn, nil
}

func (cm *dbConnectionManagerRepository) GetDBType(sourceID uuid.UUID) (string, error) {
	cm.mu.RLock()
	defer cm.mu.RUnlock()
	if dbType, ok := cm.dbTypes[sourceID]; ok {
		return dbType, nil
	}
	return "", fmt.Errorf("no database type found for sourceID: %s", sourceID)
}

func (cm *dbConnectionManagerRepository) Close(sourceID uuid.UUID) {
	cm.mu.Lock()
	defer cm.mu.Unlock()
	if c, ok := cm.connections[sourceID]; ok && c != nil {
		_ = c.GetClient().Close()
		delete(cm.connections, sourceID)
	}
	delete(cm.dbTypes, sourceID)
}

func (cm *dbConnectionManagerRepository) CloseAll() {
	cm.mu.Lock()
	defer cm.mu.Unlock()
	for id, client := range cm.connections {
		if client != nil {
			_ = client.GetClient().Close()
		}
		delete(cm.connections, id)
		delete(cm.dbTypes, id)
	}
}
