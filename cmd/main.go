package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/43nvy/wb_l0"
	"github.com/43nvy/wb_l0/internal/handler"
	"github.com/43nvy/wb_l0/internal/repository"
	"github.com/43nvy/wb_l0/internal/service"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func main() {
	logrus.SetFormatter(new(logrus.JSONFormatter))

	if err := InitConfig(); err != nil {
		logrus.Fatalf("error initializing config: %s", err.Error())
	}

	if err := godotenv.Load(); err != nil {
		logrus.Fatalf("error loading dotenv vatiables: %s", err.Error())
	}

	db, err := repository.NewPostgresDB(repository.Config{
		Host:     viper.GetString("db.host"),
		Port:     viper.GetString("db.port"),
		Username: viper.GetString("db.username"),
		DBName:   viper.GetString("db.dbname"),
		SSLMode:  viper.GetString("db.sslmode"),
		Password: os.Getenv("DB_PASSWORD"),
	})
	if err != nil {
		logrus.Fatalf("failed to initialize db: %s ", err.Error())
	}

	sc, err := service.NewNATSConn(service.Config{
		URL:     viper.GetString("nats.url"),
		Cluster: viper.GetString("nats.cluster"),
		Client:  viper.GetString("nats.client"),
	})
	if err != nil {
		logrus.Fatalf("failed to initialize nats-streaming: %s ", err.Error())
	}

	repos := repository.NewRepository(db)
	services := service.NewService(repos, sc)
	handlers := handler.NewHandler(services)

	srv := new(wb_l0.Server)
	// The following code is difficult to explain, it was taken from the Internet
	go func() {
		if err := srv.Run(viper.GetString("port"), handlers.InitRoutes()); err != nil {
			logrus.Fatalf("error occured while running http server: %s", err.Error())
		}
	}()

	logrus.Print("App started")

	handlers.InitCache()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)

	<-quit

	logrus.Print("App Shutdown")

	if err := srv.Shutdown(context.Background()); err != nil {
		logrus.Errorf("error occured on server shutting down: %s", err.Error())
	}

	if err := db.Close(); err != nil {
		logrus.Errorf("error occured on db connection close: %s", err.Error())
	}

	if err := sc.Close(); err != nil {
		logrus.Errorf("error occured on nats-streaming connection close: %s", err.Error())
	}
}

func InitConfig() error {
	viper.AddConfigPath("config")
	viper.SetConfigName("config")

	return viper.ReadInConfig()
}
