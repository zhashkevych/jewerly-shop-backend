package main

import (
	"context"
	"fmt"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"github.com/zhashkevych/jewelry-shop-backend/config"
	"github.com/zhashkevych/jewelry-shop-backend/server"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func init() {
	if err := config.Init(); err != nil {
		logrus.Fatalf("error loading config: %s\n", err.Error())
	}
}

func main() {
	db := initDB()

	srv := server.NewServer(db)

	go func() {
		if err := srv.Run(viper.GetString("port")); err != nil {
			logrus.Errorf("Error occurred while running server: %s\n", err.Error())
		}
	}()

	// graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)

	<-quit

	ctx, shutdown := context.WithTimeout(context.Background(), 5*time.Second)
	defer shutdown()

	srv.Stop(ctx)
}

func initDB() *sqlx.DB {
	db, err := sqlx.Connect("postgres", fmt.Sprintf("host=%s port=%s user=%s dbname=%s sslmode=%s password=%s",
		viper.GetString("db.postgres.host"),
		viper.GetString("db.postgres.port"),
		viper.GetString("db.postgres.username"),
		viper.GetString("db.postgres.dbname"),
		viper.GetString("db.postgres.sslmode"),
		viper.GetString("db.postgres.password"),
	))
	if err != nil {
		logrus.Fatalf("Error occurred while establishing connection to postgres: %s\n", err.Error())
	}

	if err := db.Ping(); err != nil {
		logrus.Fatalf("Error occurred while postgres health check %s\n", err.Error())
	}

	return db
}