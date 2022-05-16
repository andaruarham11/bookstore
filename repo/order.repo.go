package repo

import (
	"context"
	"errors"

	"github.com/agustadewa/book-system/configs"
	"github.com/agustadewa/book-system/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var ErrOrderNotFound = errors.New("order not found")

type Order struct {
	coll *mongo.Collection
}

func NewOrder(client *mongo.Client) *Order {
	return &Order{coll: client.Database(configs.OrderDBName).Collection(configs.OrderCollName)}
}

// Get returns an order by given order id
func (o *Order) Get(ctx context.Context, bookId string) (*models.Order, error) {
	var order models.Order
	if err := o.coll.FindOne(ctx, bson.M{"_id": bookId}).Decode(&order); err == mongo.ErrNoDocuments {
		return nil, ErrBookNotFound
	} else {
		return &order, nil
	}
}

// GetAll returns an orders by given order limit
func (o *Order) GetAll(ctx context.Context, limit int64) (*[]models.Order, error) {
	var orders []models.Order

	fr, err := o.coll.Find(ctx, bson.M{}, &options.FindOptions{Limit: &limit})
	if err = fr.All(ctx, &orders); err != nil {
		return nil, err
	}

	return &orders, nil
}

// GetAllUserId returns orders by given user id
func (o *Order) GetAllUserId(ctx context.Context, userId string) (*[]models.Order, error) {
	var orders []models.Order
	fr, err := o.coll.Find(ctx, bson.M{"user_id": userId})
	if err = fr.All(ctx, &orders); err != nil {
		return nil, err
	}
	return &orders, nil
}

// GetAllByStatus returns orders by given status
func (o *Order) GetAllByStatus(ctx context.Context, status models.OrderStatus) (*[]models.Order, error) {
	var orders []models.Order
	fr, err := o.coll.Find(ctx, bson.M{"status": status.String()})
	if err = fr.All(ctx, &orders); err != nil {
		return nil, err
	}
	return &orders, nil
}

// Add creates a new order
func (o *Order) Add(ctx context.Context, payload models.Order) (string, error) {
	if _, err := o.coll.InsertOne(ctx, payload); err != nil {
		return "", err
	}

	return payload.Id, nil
}

// UpdateStatus updates order status by given order id
func (o *Order) UpdateStatus(ctx context.Context, updatePayload models.UpdateStatusOrder) error {
	ur, err := o.coll.UpdateByID(ctx, updatePayload.OrderId, bson.M{"status": updatePayload.Status.String()})
	if err != nil {
		return err
	}
	if ur.MatchedCount == 0 {
		return ErrOrderNotFound
	}
	return nil
}
