package router

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/sing3demons/go-mongo-api/user"
	"github.com/sing3demons/go-mongo-api/utils"
	"go.mongodb.org/mongo-driver/mongo"
)

func NewRouter() *gin.Engine {
	r := gin.Default()
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

func InitUserRoutes(router *gin.Engine, db *mongo.Database) {

	r := router.Group("/api/v1")

	repository := user.NewRepository(db)
	service := user.NewService(repository)
	handler := user.NewHandler(service)
	auth := utils.Authorization()

	r.POST("/auth/register", handler.Register)
	r.POST("/auth/login", handler.Login)

	{
		r.GET("/users", auth, handler.FindAll)
		r.GET("/profile", auth, handler.Profile)
	}
}
