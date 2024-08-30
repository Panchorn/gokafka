package main

import (
	"context"
	"fmt"
	"github.com/IBM/sarama"
	"github.com/go-redis/redis/v8"
	"github.com/spf13/viper"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"logs"
	"streamA/repositories"
	"streamA/services"
	"strings"
)

func init() {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	viper.AddConfigPath("/app/config")
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	if err := viper.ReadInConfig(); err != nil {
		panic(err)
	}
}

func initDatabase() *gorm.DB {
	dsn := fmt.Sprintf("%v:%v@tcp(%v:%v)/%v?parseTime=true",
		viper.GetString("db.username"),
		viper.GetString("db.password"),
		viper.GetString("db.host"),
		viper.GetInt("db.port"),
		viper.GetString("db.database"),
	)

	dial := mysql.Open(dsn)

	db, err := gorm.Open(dial, &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		panic(err)
	}
	return db
}

func initRedis() *redis.Client {
	address := fmt.Sprintf("%v:%v", viper.GetString("redis.host"), viper.GetInt("redis.port"))
	return redis.NewClient(&redis.Options{
		Addr: address,
	})
}

func main() {
	consumer, err := sarama.NewConsumerGroup(viper.GetStringSlice("kafka.servers"), viper.GetString("kafka.group"), nil)
	if err != nil {
		panic(err)
	}
	defer consumer.Close()

	db := initDatabase()
	redisClient := initRedis()

	transactionRepository := repositories.NewTransactionRepository(db)
	transactionServiceRedis := services.NewTransactionServiceRedis(transactionRepository, redisClient)
	transferEventHandler := services.NewTransferEventHandler(transactionRepository, transactionServiceRedis)
	transferConsumerHandler := services.NewConsumerHandler(transferEventHandler)

	logs.Info("streamA started...")
	for {
		consumer.Consume(context.Background(), viper.GetStringSlice("kafka.topic-subscriptions"), transferConsumerHandler)
	}
}
