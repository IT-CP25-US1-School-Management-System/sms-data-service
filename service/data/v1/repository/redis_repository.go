package repository

import (
	"github.com/GodeFvt/go-backend/redis"
	"github.com/IT-CP25-US1-School-Management-System/sms-data-service/service/data/v1"
)

type redisRepository struct {
	client *redis.Client
}

func NewRedisRepository(client *redis.Client) data.RedisRepository {
	return &redisRepository{
		client: client,
	}
}
