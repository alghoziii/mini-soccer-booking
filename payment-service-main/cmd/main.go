package cmd

import (
	"context"
	"fmt"
	"github.com/didip/tollbooth"
	"github.com/didip/tollbooth/limiter"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/spf13/cobra"
	"net/http"
	"payment-service/clients"
	midtransClient "payment-service/clients/midtrans"
	"payment-service/common/response"
	"payment-service/common/storage"
	"payment-service/config"
	"payment-service/constants"
	"payment-service/controllers/http"
	kafkaClient "payment-service/controllers/kafka"
	"payment-service/domain/models"
	"payment-service/middlewares"
	"payment-service/repositories"
	"payment-service/routes"
	"payment-service/services"
	"time"
)

var command = &cobra.Command{
	Use:   "serve",
	Short: "Start the server",
	Run: func(c *cobra.Command, args []string) {
		_ = godotenv.Load()
		config.Init()
		db, err := config.InitDatabase()
		if err != nil {
			panic(err)
		}

		loc, err := time.LoadLocation("Asia/Jakarta")
		if err != nil {
			panic(err)
		}
		time.Local = loc

		err = db.AutoMigrate(
			&models.Payment{},
			&models.PaymentHistory{},
		)
		if err != nil {
			panic(err)
		}

		storageClient := initS3()
		kafka := kafkaClient.NewKafkaRegistry(config.Config.Kafka.Brokers)
		midtrans := midtransClient.NewMidtransClient(
			config.Config.Midtrans.ServerKey,
			config.Config.Midtrans.IsProduction,
		)
		client := clients.NewClientRegistry()
		repository := repositories.NewRepositoryRegistry(db)
		service := services.NewServiceRegistry(repository, storageClient, kafka, midtrans)
		controller := controllers.NewControllerRegistry(service)

		router := gin.Default()
		router.Use(middlewares.HandlePanic())
		router.NoRoute(func(c *gin.Context) {
			c.JSON(http.StatusNotFound, response.Response{
				Status:  constants.Error,
				Message: fmt.Sprintf("Path %s", http.StatusText(http.StatusNotFound)),
			})
		})
		router.GET("/", func(c *gin.Context) {
			c.JSON(http.StatusOK, response.Response{
				Status:  constants.Success,
				Message: "Welcome to Payment Service",
			})
		})
		router.Use(func(c *gin.Context) {
			c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
			c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, PATCH")
			c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, x-service-name, x-request-at, x-api-key")
			if c.Request.Method == "OPTIONS" {
				c.AbortWithStatus(204)
				return
			}
			c.Next()
		})

		lmt := tollbooth.NewLimiter(
			config.Config.RateLimiterMaxRequest,
			&limiter.ExpirableOptions{
				DefaultExpirationTTL: time.Duration(config.Config.RateLimiterTimeSecond) * time.Second,
			})
		router.Use(middlewares.RateLimiter(lmt))

		group := router.Group("/api/v1")
		route := routes.NewRouteRegistry(controller, group, client)
		route.Serve()

		port := fmt.Sprintf(":%d", config.Config.Port)
		router.Run(port)
	},
}

func Run() {
	err := command.Execute()
	if err != nil {
		panic(err)
	}
}

func initS3() storage.IClient {
	storageClient, err := storage.NewS3Client(context.Background(), storage.S3Config{
		Region:          config.Config.S3.Region,
		BucketName:      config.Config.S3.BucketName,
		BaseURL:         config.Config.S3.BaseURL,
		AccessKeyID:     config.Config.S3.AccessKeyID,
		SecretAccessKey: config.Config.S3.SecretAccessKey,
	})
	if err != nil {
		panic(err)
	}

	return storageClient
}
