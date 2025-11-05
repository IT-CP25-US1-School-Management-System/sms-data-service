package repository

import (
	"context"
	"fmt"
	"net/url"
	"sync"
	"time"

	"github.com/GodeFvt/go-backend/psql"
	"github.com/IT-CP25-US1-School-Management-System/sms-data-service/models/entity"
	"github.com/IT-CP25-US1-School-Management-System/sms-data-service/service/data/v1"
	"github.com/IT-CP25-US1-School-Management-System/sms-data-service/service/database/v1"
	"github.com/IT-CP25-US1-School-Management-System/sms-data-service/utils/crypto"
	"github.com/gofrs/uuid"
)

type dbConnectionManagerRepository struct {
	datasetRepo  data.PsqlDatasetRepository
	connections  map[uuid.UUID]*psql.Client
	mu           sync.RWMutex
	cryptoSecret string
}

func NewDBConnectionManagerRepository(datasetRepo data.PsqlDatasetRepository, cryptoSecret string) database.DBConnectionManagerRepository {
	return &dbConnectionManagerRepository{
		datasetRepo:  datasetRepo,
		cryptoSecret: cryptoSecret,
		connections:  make(map[uuid.UUID]*psql.Client),
	}
}

func (cm *dbConnectionManagerRepository) createConnectionDetails(source *entity.Sources) (string, psql.Driver, error) {
	decryptedPass, err := crypto.Decrypt(source.Password, cm.cryptoSecret)
	if err != nil {
		return "", psql.Postgres, fmt.Errorf("failed to decrypt password: %w", err)
	}

	switch source.DBType {
	case "postgres":
		user := url.QueryEscape(source.Username)
		pass := url.QueryEscape(decryptedPass)
		host := source.Host
		connStr := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable",
			user, pass, host, source.Port, source.DatabaseName)
		return connStr, psql.Postgres, nil

	case "mysql":
		user := url.QueryEscape(source.Username)
		pass := url.QueryEscape(decryptedPass)
		dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?parseTime=true&charset=utf8mb4&loc=Local",
			user, pass, source.Host, source.Port, source.DatabaseName)
		return dsn, psql.MySQL, nil

	default:
		return "", psql.Postgres, fmt.Errorf("unsupported database type: %s", source.DBType)
	}
}

func (cm *dbConnectionManagerRepository) GetConnection(ctx context.Context, sourceID uuid.UUID) (*psql.Client, error) {
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

	connStr, drv, err := cm.createConnectionDetails(source)
	if err != nil {
		return nil, err
	}

	// สร้าง connection
	client, err := psql.NewConnection(connStr, drv)
	if err != nil {
		return nil, fmt.Errorf("failed to create connection: %w", err)
	}

	// ใส่ cache โดยตรวจซ้ำ
	cm.mu.Lock()
	if existing, ok := cm.connections[sourceID]; ok {
		cm.mu.Unlock()
		// มีคนสร้างทันก่อน ใช้อันเดิมและปิดอันใหม่
		client.GetClient().Close()
		return existing, nil
	}
	cm.connections[sourceID] = client
	cm.mu.Unlock()
	return client, nil
}

func (cm *dbConnectionManagerRepository) Close(sourceID uuid.UUID) {
	cm.mu.Lock()
	defer cm.mu.Unlock()
	if c, ok := cm.connections[sourceID]; ok && c != nil {
		_ = c.GetClient().Close()
		delete(cm.connections, sourceID)
	}
}

func (cm *dbConnectionManagerRepository) CloseAll() {
	cm.mu.Lock()
	defer cm.mu.Unlock()
	for id, client := range cm.connections {
		if client != nil {
			_ = client.GetClient().Close()
		}
		delete(cm.connections, id)
	}
}
