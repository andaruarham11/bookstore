package handlers

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/agustadewa/book-system/models"
	"github.com/agustadewa/book-system/repo"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func NewOrder(engine *gin.Engine, client *mongo.Client) *OrderHandler {
	return &OrderHandler{
		engine: engine,
		order:  repo.NewOrder(client),
	}
}

type OrderHandler struct {
	engine *gin.Engine
	order  *repo.Order
}

func (h *OrderHandler) RegisterEndpoint() {
	h.engine.POST("/order", h.addOrder)
	h.engine.GET("/order/:order_id", h.getOrder)
	h.engine.GET("/order/all", h.getAllOrders)
	h.engine.GET("/order/all/byuserid/:user_id", h.getAllOrdersByUserId)
	h.engine.GET("/order/all/bystatus/:status", h.getAllOrdersByStatus)
	h.engine.PUT("/order/:order_id/setstatus/:status", h.setOrderStatus)
}

func (h *OrderHandler) addOrder(c *gin.Context) {
	ctx := c.Request.Context()

	var addOrder models.AddOrder
	if err := c.BindJSON(&addOrder); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("error parsing request data: %s", err)})
		return
	}

	// validate order status
	orderStatus, err := models.IsValidOrderStatus(addOrder.Status.String())
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// get order
	_, err = h.order.GetByBookId(ctx, addOrder.BookId)
	if err == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": repo.ErrOrderExists.Error()})
		return
	}
	if err != nil && err != repo.ErrOrderNotFound {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}

	// add order
	addOrderPayload := models.Order{
		Id:        primitive.NewObjectID().Hex(),
		UserId:    addOrder.UserId,
		BookId:    addOrder.BookId,
		Qty:       addOrder.Qty,
		OrderTime: addOrder.OrderTime,
		Status:    orderStatus,
	}
	id, err := h.order.Add(ctx, addOrderPayload)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"result":  gin.H{"id": id},
	})
}

func (h *OrderHandler) getOrder(c *gin.Context) {
	ctx := c.Request.Context()

	orderId := c.Param("order_id")

	// get order
	order, err := h.order.Get(ctx, orderId)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"result":  order,
	})
}

func (h *OrderHandler) getAllOrders(c *gin.Context) {
	ctx := c.Request.Context()

	limitStr := c.Request.URL.Query().Get("limit")
	limit, _ := strconv.ParseInt(limitStr, 10, 64)

	if limit < 10 || limit > 100 {
		limit = 10
	}
	orders, err := h.order.GetAll(ctx, limit)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "result": orders})
}

func (h *OrderHandler) getAllOrdersByUserId(c *gin.Context) {
	ctx := c.Request.Context()

	userId := c.Param("user_id")

	limitStr := c.Request.URL.Query().Get("limit")
	limit, _ := strconv.ParseInt(limitStr, 10, 64)

	if limit < 10 || limit > 100 {
		limit = 10
	}
	orders, err := h.order.GetAllByUserId(ctx, userId, limit)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "result": orders})
}

func (h *OrderHandler) getAllOrdersByStatus(c *gin.Context) {
	ctx := c.Request.Context()

	// validate order status
	status := c.Param("status")
	orderStatus, err := models.IsValidOrderStatus(status)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// check limit
	limitStr := c.Request.URL.Query().Get("limit")
	limit, _ := strconv.ParseInt(limitStr, 10, 64)
	if limit < 10 || limit > 100 {
		limit = 10
	}

	orders, err := h.order.GetAllByUserId(ctx, orderStatus.String(), limit)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "result": orders})
}

func (h *OrderHandler) setOrderStatus(c *gin.Context) {
	ctx := c.Request.Context()

	orderId := c.Param("order_id")

	// validate order status
	status := c.Param("status")
	orderStatus, err := models.IsValidOrderStatus(status)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err = h.order.UpdateStatus(ctx, orderId, orderStatus); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"result":  fmt.Sprintf("order status setted to %s", orderStatus.String()),
	})
}
