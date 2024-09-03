package services

import (
	"apiA/repositories"
	"context"
	"encoding/json"
	"github.com/go-redis/redis/v8"
	"github.com/spf13/viper"
	"logs"
	"time"
)

type transactionServiceRedis struct {
	transactionRepository repositories.TransactionRepository
	redisClient           *redis.Client
}

type TransactionService interface {
	GetTransaction(requestID string, refID string) (repositories.Transaction, error)
}

func NewTransactionServiceRedis(transactionRepository repositories.TransactionRepository, redisClient *redis.Client) TransactionService {
	return transactionServiceRedis{transactionRepository, redisClient}
}

func (r transactionServiceRedis) GetTransaction(requestID string, refID string) (transaction repositories.Transaction, err error) {
	logs.Info(requestID, "getting transaction from redis")
	key := "service:transactions:" + refID

	// Redis Get
	transactionJson, err := r.redisClient.Get(context.Background(), key).Result()
	if err == nil {
		err = json.Unmarshal([]byte(transactionJson), &transaction)
		if err == nil {
			logs.Info(requestID, "transaction cache "+transaction.ToString())
			return transaction, nil
		}
	}

	// Database Get
	transaction, err = r.transactionRepository.FindByRefID(refID)
	if err != nil {
		return repositories.Transaction{}, err
	}

	// Redis Set
	data, err := json.Marshal(transaction)
	if err != nil {
		return repositories.Transaction{}, err
	}

	ttl := viper.GetDuration("redis.ttl.get_transaction_in_ms")
	err = r.redisClient.Set(context.Background(), key, string(data), time.Millisecond*ttl).Err()
	if err != nil {
		return repositories.Transaction{}, err
	}

	logs.Info(requestID, "transaction cache "+transaction.ToString())
	return transaction, nil
}
