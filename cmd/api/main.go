package main

import (
	"context"
	_ "github.com/lib/pq"
	"github.com/minio/minio-go"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"github.com/zhashkevych/jewelry-shop-backend"
	"github.com/zhashkevych/jewelry-shop-backend/pkg/config"
	"github.com/zhashkevych/jewelry-shop-backend/pkg/email"
	"github.com/zhashkevych/jewelry-shop-backend/pkg/handler"
	"github.com/zhashkevych/jewelry-shop-backend/pkg/payment"
	"github.com/zhashkevych/jewelry-shop-backend/pkg/repository"
	"github.com/zhashkevych/jewelry-shop-backend/pkg/service"
	"github.com/zhashkevych/jewelry-shop-backend/pkg/storage"
	"io"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func init() {
	logrus.SetLevel(logrus.DebugLevel)
	logrus.SetFormatter(&logrus.JSONFormatter{})

	file, err := os.OpenFile("/logs/api.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		logrus.Fatal(err)
	}

	mw := io.MultiWriter(os.Stdout, file)
	logrus.SetOutput(mw)

	if err := config.Init(); err != nil {
		logrus.Fatalf("error loading config: %s\n", err.Error())
	}
}

func main() {
	// Init infrastructure layer
	db, err := repository.NewPostgresDB(repository.Config{
		Host:     viper.GetString("db.postgres.host"),
		Port:     viper.GetString("db.postgres.port"),
		Username: viper.GetString("db.postgres.username"),
		DBName:   viper.GetString("db.postgres.dbname"),
		SSLMode:  viper.GetString("db.postgres.sslmode"),
		Password: os.Getenv("POSTGRES_PASSWORD"),
	})
	if err != nil {
		logrus.Fatalf("Error occurred on db initialization: %s\n", err.Error())
	}

	minioStorage, err := initStorage()
	if err != nil {
		logrus.Fatalf("Error occurred on storage initialization: %s\n", err.Error())
	}

	apiKey := os.Getenv("PAYMENT_API_KEY")
	if apiKey == "" {
		logrus.Fatalln("Payment credentials are empty")
	}

	paymentProvider := payment.NewIsracardProvider(
		viper.GetString("payments.endpoint"), apiKey,
		viper.GetString("payments.return_url"), viper.GetString("payments.callback_url"))

	emailPassword := os.Getenv("EMAIL_PASSWORD")
	if apiKey == "" {
		logrus.Fatalln("Email password is empty")
	}

	emailSender := email.NewSMTPClient(
		viper.GetString("email.smtp.host"),
		viper.GetString("email.smtp.port"),
		viper.GetString("email.sender.email"),
		emailPassword)

	// Init Dependecies
	repos := repository.NewRepository(db)
	services := service.NewServices(service.Dependencies{
		Repos:           repos,
		HashSalt:        viper.GetString("auth.hash_salt"),
		SigningKey:      []byte(viper.GetString("auth.signing_key")),
		FileStorage:     minioStorage,
		PaymentProvider: paymentProvider,
		SupportEmail:    viper.GetString("email.support.email"),
		SupportName:     viper.GetString("email.support.name"),
		SenderName:      viper.GetString("email.sender.name"),
		SenderEmail:     viper.GetString("email.sender.email"),

		OrderInfoSupportTemplate: viper.GetString("email.templates.order_info_support"),
		OrderInfoSupportSubject:  viper.GetString("email.subjects.order_info_support"),

		OrderInfoCustomerTemplate: viper.GetString("email.templates.order_info_customer"),
		OrderInfoCustomerSubject:  viper.GetString("email.subjects.order_info_customer"),

		PaymentInfoSupportTemplate: viper.GetString("email.templates.payment_info_support"),
		PaymentInfoSupportSubject:  viper.GetString("email.subjects.payment_info_support"),

		PaymentInfoCustomerTemplate: viper.GetString("email.templates.payment_info_customer"),
		PaymentInfoCustomerSubject:  viper.GetString("email.subjects.payment_info_customer"),

		EmailSender: emailSender,
	})
	handlers := handler.NewHandler(services)

	// Create & Run HTTP Server
	server := jewerly.NewServer()
	go func() {
		if err := server.Run(viper.GetString("port"), handlers.Init()); err != nil {
			logrus.Errorf("Error occurred while running server: %s\n", err.Error())
		}
	}()

	logrus.Info("Application Started")

	// graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)

	<-quit

	ctx, shutdown := context.WithTimeout(context.Background(), 5*time.Second)
	defer shutdown()

	if err := server.Stop(ctx); err != nil {
		logrus.Errorf("error occurred while shutting down http server: %s\n", err.Error())
	}

	if err := db.Close(); err != nil {
		logrus.Errorf("error occurred while closing db connection: %s\n", err.Error())
	}
}

func initStorage() (storage.Storage, error) {
	client, err := minio.New(
		viper.GetString("storage.url"),
		os.Getenv("ACCESS_KEY"),
		os.Getenv("SECRET_KEY"), false)
	if err != nil {
		return nil, err
	}

	exists, err := client.BucketExists(viper.GetString("storage.bucket"))
	if err != nil {
		return nil, err
	}

	logrus.Infof("Bucket %s exists: %v", viper.GetString("storage.bucket"), exists)

	return storage.NewFileStorage(client,
		viper.GetString("storage.bucket"),
		viper.GetString("storage.url"),
		os.Getenv("HOST")), nil
}
