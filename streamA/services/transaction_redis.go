package services

import (
	"context"
	"github.com/go-redis/redis/v8"
	"streamA/repositories"
)

type transactionServiceRedis struct {
	transactionRepository repositories.TransactionRepository
	redisClient           *redis.Client
}

type TransactionService interface {
	EvictTransaction(refID string) error
}

func NewTransactionServiceRedis(transactionRepository repositories.TransactionRepository, redisClient *redis.Client) TransactionService {
	return transactionServiceRedis{transactionRepository, redisClient}
}

func (r transactionServiceRedis) EvictTransaction(refID string) error {
	key := "service:transactions:" + refID
	_, err := r.redisClient.Del(context.Background(), key).Result()
	if err != nil {
		return err
	}
	return nil
}
