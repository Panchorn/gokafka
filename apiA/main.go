package main

import (
	"apiA/controllers"
	"apiA/services"
	"github.com/IBM/sarama"
	"github.com/gofiber/fiber/v3"
	"github.com/spf13/viper"
	"log"
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
	producer, err := sarama.NewSyncProducer(viper.GetStringSlice("kafka.servers"), nil)
	if err != nil {
		panic(err)
	}
	defer producer.Close()

	eventProducer := services.NewEventProducer(producer)
	transferService := services.NewTransferService(eventProducer)
	transferController := controllers.NewTransferController(transferService)

	app := fiber.New()

	app.Post("/transfer", transferController.Transfer)

	log.Println("Starting app on port 8000")
	app.Listen(":8000")
}
