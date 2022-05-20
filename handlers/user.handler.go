package handlers

import (
	"encoding/base64"
	"fmt"
	"net/http"

	"github.com/agustadewa/book-system/models"
	"github.com/agustadewa/book-system/repo"
	"github.com/gin-gonic/gin"
	jsoniter "github.com/json-iterator/go"
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

func (h *UserHandler) RegisterEndpoints() {
	h.engine.POST("/login", h.login)
	h.engine.GET("/logout", h.logout)
	h.engine.POST("/register", h.register)
	h.engine.DELETE("/user/:user_id", h.delete)
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

	var userResp models.UserResp
	bytes, _ := jsoniter.Marshal(user)
	_ = jsoniter.Unmarshal(bytes, &userResp)

	c.JSON(http.StatusOK, gin.H{"success": true, "result": userResp})
}

func (h *UserHandler) logout(c *gin.Context) {
	c.SetCookie("authenticated", "", -1, "/", "localhost", false, true)

	c.JSON(http.StatusOK, gin.H{"success": true, "result": "successfully logged out"})
}

func (h *UserHandler) register(c *gin.Context) {
	ctx := c.Request.Context()
	var register models.AddUser
	if err := c.BindJSON(&register); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("error parsing request data: %s", err)})
		return
	}

	// Check existing user
	_, err := h.user.GetByUserName(ctx, register.Username)
	if err == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": repo.ErrUserExists.Error()})
		return
	}
	if err != nil && err != repo.ErrUserNotFound {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Add user
	addPayload := models.User{
		Id:         primitive.NewObjectID().Hex(),
		Name:       register.Name,
		UserName:   register.Username,
		Password:   base64.StdEncoding.EncodeToString([]byte(register.Password)),
		IsAdmin:    register.IsAdmin,
		IsVerified: false,
	}

	id, err := h.user.Add(ctx, addPayload)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "result": gin.H{"id": id}})
}

func (h *UserHandler) delete(c *gin.Context) {
	ctx := c.Request.Context()

	userId := c.Param("user_id")

	if err := h.user.Delete(ctx, userId); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "result": fmt.Sprintf("user %s has been deleted", userId)})
}
