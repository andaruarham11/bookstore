package repo

import (
	"context"
	"errors"

	"github.com/agustadewa/book-system/configs"
	"github.com/agustadewa/book-system/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

var ErrPaymentNotFound = errors.New("payment not found")
var ErrPaymentExists = errors.New("payment already exists")

type Payment struct {
	coll *mongo.Collection
}

func NewPayment(client *mongo.Client) *Payment {
	return &Payment{coll: client.Database(configs.PaymentDBName).Collection(configs.PaymentCollName)}
}

// GetByUserId returns a payment by given user id
func (p *Payment) GetByUserId(ctx context.Context, userId string) (*models.Payment, error) {
	var payment models.Payment
	if err := p.coll.FindOne(ctx, bson.M{"user_id": userId}).Decode(&payment); err == mongo.ErrNoDocuments {
		return nil, ErrPaymentNotFound
	} else {
		return &payment, nil
	}
}

// GetByOrderId returns a payment by given order id
func (p *Payment) GetByOrderId(ctx context.Context, orderId string) (*models.Payment, error) {
	var payment models.Payment
	if err := p.coll.FindOne(ctx, bson.M{"order_id": orderId}).Decode(&payment); err == mongo.ErrNoDocuments {
		return nil, ErrPaymentNotFound
	} else {
		return &payment, nil
	}
}

// Add creates a new payment
func (p *Payment) Add(ctx context.Context, payload models.Payment) (string, error) {
	if _, err := p.coll.InsertOne(ctx, payload); err != nil {
		return "", err
	}

	return payload.Id, nil
}
