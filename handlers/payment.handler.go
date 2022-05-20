package handlers

import (
	"fmt"
	"net/http"

	"github.com/agustadewa/book-system/models"
	"github.com/agustadewa/book-system/repo"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func NewPayment(engine *gin.Engine, client *mongo.Client) *PaymentHandler {
	return &PaymentHandler{
		engine:  engine,
		payment: repo.NewPayment(client),
		user:    repo.NewUser(client),
		order:   repo.NewOrder(client),
	}
}

type PaymentHandler struct {
	engine  *gin.Engine
	payment *repo.Payment
	user    *repo.User
	order   *repo.Order
}

func (h *PaymentHandler) RegisterEndpoints() {
	h.engine.POST("/payment", h.addPayment)
	h.engine.GET("/payment/:payment_id", h.getPayment)
	h.engine.GET("/payment/byorderid/:order_id", h.getPaymentByOrderId)
	h.engine.DELETE("/payment/:payment_id", h.delete)
}

func (h *PaymentHandler) addPayment(c *gin.Context) {
	ctx := c.Request.Context()

	var addPayment models.AddPayment
	if err := c.BindJSON(&addPayment); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("error parsing request data: %s", err)})
		return
	}

	// check existing payment
	_, err := h.payment.GetByOrderId(ctx, addPayment.OrderId)
	if err == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": repo.ErrPaymentExists.Error()})
		return
	}
	if err != nil && err != repo.ErrPaymentNotFound {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// check existing order
	if _, err = h.order.Get(ctx, addPayment.OrderId); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// check existing user
	if _, err = h.user.Get(ctx, addPayment.UserId); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// add payment
	addPaymentPayload := models.Payment{
		Id:      primitive.NewObjectID().Hex(),
		UserId:  addPayment.UserId,
		OrderId: addPayment.OrderId,
		Receipt: addPayment.Receipt,
	}
	id, err := h.payment.Add(ctx, addPaymentPayload)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// update order status
	if err = h.order.UpdateStatus(ctx, addPayment.OrderId, models.Paid); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "result": gin.H{"id": id}})
}

func (h *PaymentHandler) getPaymentByOrderId(c *gin.Context) {
	ctx := c.Request.Context()

	orderId := c.Param("order_id")

	payment, err := h.payment.GetByOrderId(ctx, orderId)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if payment == nil {
		payment = &models.Payment{}
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "result": payment})
}

func (h *PaymentHandler) getPaymentByUserId(c *gin.Context) {
	ctx := c.Request.Context()

	userId := c.Param("user_id")

	payment, err := h.payment.GetByUserId(ctx, userId)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if payment == nil {
		payment = &models.Payment{}
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "result": payment})
}

func (h *PaymentHandler) getPayment(c *gin.Context) {
	ctx := c.Request.Context()

	paymentId := c.Param("payment_id")

	payment, err := h.payment.Get(ctx, paymentId)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if payment == nil {
		payment = &models.Payment{}
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "result": payment})
}

func (h *PaymentHandler) delete(c *gin.Context) {
	ctx := c.Request.Context()

	paymentId := c.Param("payment_id")

	if err := h.payment.Delete(ctx, paymentId); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "result": fmt.Sprintf("payment %s has been deleted", paymentId)})
}
