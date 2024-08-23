package main

import (
	"context"
	"fmt"
	"github.com/IBM/sarama"
	"github.com/spf13/viper"
	"streamB/services"
	"strings"
)

func init() {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	if err := viper.ReadInConfig(); err != nil {
		panic(err)
	}
}

func main() {
	consumer, err := sarama.NewConsumerGroup(viper.GetStringSlice("kafka.servers"), viper.GetString("kafka.group"), nil)
	if err != nil {
		panic(err)
	}
	defer consumer.Close()

	eventHandler := services.NewEventHandler()
	consumerHandler := services.NewConsumerHandler(eventHandler)

	fmt.Println("streamB started...")
	topic := viper.GetStringSlice("kafka.topic-subscriptions")
	for {
		consumer.Consume(context.Background(), topic, consumerHandler)
	}
}
