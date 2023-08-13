package user

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

type handler struct {
	service *service
}

type MyContext interface {
	GetUserId() (string, error)
	SetUserId(id string)
	BindJSON(obj interface{}) error
	JSON(code int, obj interface{})
}

func NewHandler(service *service) *handler {
	return &handler{service}
}

func (u *handler) Register(c *gin.Context) {
	var user User
	err := c.BindJSON(&user)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	uu, err := u.service.Register(user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"user": uu})
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
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	user, err := u.service.FindOne(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"user": user})
}

func (u *handler) FindAll(c *gin.Context) {

	users, err := u.service.FindAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"users": users})
}

func (h *handler) Login(c *gin.Context) {

	var user Login
	err := c.BindJSON(&user)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	token, err := h.service.Login(user.Email, user.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": token})
}
