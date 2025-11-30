package adapters

import (
	database "github.com/IT-CP25-US1-School-Management-System/sms-data-service/service/database/v1"
)

// RegisterAllAdapters registers all built-in database adapters
// This is the single entry point for registering complete database support
func RegisterAllAdapters() error {
	adapters := []database.DatabaseAdapter{
		NewPostgresAdapter(),
		NewMySQLAdapter(),
		NewMSSQLAdapter(),
	}

	for _, adapter := range adapters {
		if err := database.RegisterAdapter(adapter); err != nil {
			return err
		}
	}

	return nil
}
