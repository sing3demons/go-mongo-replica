package user

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sing3demons/go-mongo-api/response"
	"github.com/sirupsen/logrus"
	"go.uber.org/zap"
)

type handler struct {
	service *service
	logger  *zap.Logger
}

func NewHandler(service *service, logger *zap.Logger) *handler {
	return &handler{service: service, logger: logger}
}

func (u *handler) Register(c *gin.Context) {
	var user User

	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	uu, err := u.service.Register(user)
	if err != nil {
		// c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		response.Json(c, http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	if err != nil {
		// c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		response.Json(c, http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	// c.JSON(http.StatusOK, gin.H{"user": uu})
	response.Json(c, http.StatusOK, gin.H{"user": uu})
}

func (u *handler) getUserId(c *gin.Context) (string, error) {
	id, ok := c.Get("userId")
	if !ok {
		return "", fmt.Errorf("unauthorized")
	}

	return id.(string), nil
}

func (u *handler) Profile(c *gin.Context) {
	id, err := u.getUserId(c)
	if err != nil {
		// c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		response.Json(c, http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	user, err := u.service.FindOne(id)
	if err != nil {
		// c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		response.Json(c, http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	// c.JSON(http.StatusOK, gin.H{"user": user})
	response.Json(c, http.StatusOK, gin.H{"user": user})
}

func (u *handler) FindAll(c *gin.Context) {

	users, err := u.service.FindAll()
	if err != nil {
		// c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		response.Json(c, http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// c.JSON(http.StatusOK, gin.H{"users": users})
	response.Json(c, http.StatusOK, gin.H{"users": users})
}

func (h *handler) logInfo(msg string, fields any) {
	var v map[string]any
	fieldsByte, _ := json.Marshal(fields)
	json.Unmarshal(fieldsByte, &v)

	if v["password"] != nil {
		v["password"] = "****"
	}

	if v["Password"] != nil {
		v["Password"] = "****"
	}
	h.logger.Info(msg, zap.Any("commonFields", v))
}

func (h *handler) Login(c *gin.Context) {
	var user Login
	if err := c.ShouldBindJSON(&user); err != nil {
		// c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		response.Json(c, http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	h.logInfo("login", user)

	token, err := h.service.Login(user.Email, user.Password)
	if err != nil {
		logrus.Error("error login: ", err)
		// c.JSON(http.StatusNotFound, gin.H{"error": "email or password wrong"})
		response.Json(c, http.StatusNotFound, gin.H{"error": "email or password wrong"})
		return
	}

	response.Json(c, http.StatusOK, gin.H{"token": token})
}
