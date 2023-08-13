package router

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	logMiddleware "github.com/sing3demons/go-mongo-api/logger"
	"github.com/sing3demons/go-mongo-api/user"
	"github.com/sing3demons/go-mongo-api/utils"
	prometheus "github.com/zsais/go-gin-prometheus"
	"go.mongodb.org/mongo-driver/mongo"
	"go.uber.org/zap"
)

func NewRouter(logger *zap.Logger) *gin.Engine {
	r := gin.Default()
	p := prometheus.NewPrometheus("gin")
	p.Use(r)
	r.Use(logMiddleware.ZapLogger(logger))
	r.Use(logMiddleware.RecoveryWithZap(logger, true))
	config := cors.DefaultConfig()
	config.AllowOrigins = []string{
		"http://localhost:8081",
	}
	config.AllowHeaders = []string{
		"Origin",
		"Authorization",
	}

	r.Use(cors.New(config))
	return r
}

func InitUserRoutes(router *gin.Engine, db *mongo.Database, logger *zap.Logger) {

	r := router.Group("/api/v1")

	repository := user.NewRepository(db, logger)
	service := user.NewService(repository, logger)
	handler := user.NewHandler(service, logger)
	auth := utils.Authorization()

	r.POST("/auth/register", handler.Register)
	r.POST("/auth/login", handler.Login)

	{
		r.GET("/users", auth, handler.FindAll)
		r.GET("/profile", auth, handler.Profile)
	}
}
