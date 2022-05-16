package handlers

import (
	"encoding/base64"
	"fmt"
	"net/http"

	"github.com/agustadewa/book-system/models"
	"github.com/agustadewa/book-system/repo"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func NewUser(engine *gin.Engine, client *mongo.Client) *UserHandler {
	return &UserHandler{
		engine: engine,
		user:   repo.NewUser(client),
	}
}

type UserHandler struct {
	engine *gin.Engine
	user   *repo.User
}

func (h *UserHandler) RegisterEndpoint() {
	h.engine.POST("/login", h.login)
	h.engine.GET("/logout", h.logout)
	h.engine.POST("/register", h.register)
}

func (h *UserHandler) login(c *gin.Context) {
	ctx := c.Request.Context()

	var login models.Login
	if err := c.BindJSON(&login); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("error parsing request data: %s", err)})
		return
	}

	// Get user
	user, err := h.user.GetByUserName(ctx, login.Username)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if user.Password != base64.StdEncoding.EncodeToString([]byte(login.Password)) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "password is incorrect"})
		return
	}

	c.SetCookie("authenticated", "true", 60*60*24, "/", "localhost", false, true)

	c.JSON(http.StatusOK, gin.H{"success": true})
}

func (h *UserHandler) logout(c *gin.Context) {
	c.SetCookie("authenticated", "", -1, "/", "localhost", false, true)

	c.JSON(http.StatusOK, gin.H{"success": true})
}

func (h *UserHandler) register(c *gin.Context) {
	ctx := c.Request.Context()
	var register models.AddUser
	if err := c.BindJSON(&register); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("error parsing request data: %s", err)})
		return
	}

	// Get user
	_, err := h.user.GetByUserName(ctx, register.Username)
	if err == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": repo.ErrUserExists.Error()})
		return
	}
	if err != nil && err != repo.ErrUserNotFound {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}

	// Add user
	addPayload := models.User{
		Id:         primitive.NewObjectID().Hex(),
		Name:       register.Name,
		UserName:   register.Username,
		Password:   base64.StdEncoding.EncodeToString([]byte(register.Password)),
		IsAdmin:    false,
		IsVerified: false,
	}
	if _, err = h.user.Add(ctx, addPayload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true})
}
