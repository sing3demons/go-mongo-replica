package main

import (
	"context"
	"errors"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/sing3demons/go-mongo-api/database"
	"github.com/sing3demons/go-mongo-api/logger"
	"github.com/sing3demons/go-mongo-api/router"
	"github.com/sirupsen/logrus"
)

func main() {
	if os.Getenv(gin.EnvGinMode) != gin.ReleaseMode {
		if err := godotenv.Load(".env"); err != nil {
			logrus.Fatal("Error loading .env file")
		}
	}

	uri := os.Getenv("MONGO_URL")
	if uri == "" {
		logrus.Fatal("missing connection string")
	}

	client := database.Connect(uri)
	defer func() {
		if err := client.Disconnect(context.TODO()); err != nil {
			panic(err)
		}
	}()

	db := client.Database("users")

	logger, _ := logger.NewLogger()
	defer logger.Sync()

	r := router.NewRouter(logger)
	router.InitUserRoutes(r, db, logger)
	ServeHttp(":8080", "user-service", r)

}

func ServeHttp(addr, serviceName string, router http.Handler) {
	srv := &http.Server{
		Addr:    addr,
		Handler: router,
	}

	go func() {
		logrus.Infof("[%s] http listen: %v", serviceName, srv.Addr)

		if err := srv.ListenAndServe(); err != nil && errors.Is(err, http.ErrServerClosed) {
			logrus.Error("server listen err: ", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logrus.Warn("shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		logrus.Fatal("server forced to shutdown: ", err)
	}

	logrus.Warn("server exited")
}
