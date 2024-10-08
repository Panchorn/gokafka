package main

import (
	"apiA/controllers"
	"apiA/repositories"
	"apiA/services"
	"fmt"
	"github.com/IBM/sarama"
	"github.com/go-redis/redis/v8"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/spf13/viper"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"logs"
	"strings"
	"time"
)

func init() {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	viper.AddConfigPath("/app")
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

	fmt.Println(dsn)

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
	db := initDatabase()
	redisClient := initRedis()
	transactionRepository := repositories.NewTransactionRepository(db)

	producer, err := sarama.NewSyncProducer(viper.GetStringSlice("kafka.servers"), nil)
	if err != nil {
		panic(err)
	}
	defer producer.Close()

	eventProducer := services.NewEventProducer(producer)
	transactionServiceRedis := services.NewTransactionServiceRedis(transactionRepository, redisClient)
	transferService := services.NewTransferService(eventProducer, transactionRepository, transactionServiceRedis)
	transferController := controllers.NewTransferController(transferService)

	app := echo.New()

	app.Use(middleware.TimeoutWithConfig(middleware.TimeoutConfig{
		Skipper:      middleware.DefaultSkipper,
		ErrorMessage: "timeout during request handling",
		OnTimeoutRouteErrorHandler: func(err error, c echo.Context) {
			logs.Error("main", "timeout to call "+c.Path())
		},
		Timeout: viper.GetDuration("app.http_timeout_in_ms") * time.Millisecond,
	}))
	app.Use(middleware.RequestID())

	app.POST("/transfers", transferController.Transfer)
	app.GET("/transfers/transactions", transferController.TransferTransactions)

	app.Logger.Fatal(app.Start(fmt.Sprintf(":%v", viper.GetInt("app.port"))))
}
