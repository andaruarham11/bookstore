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

var ErrUserNotFound = errors.New("user not found")
var ErrUserExists = errors.New("user already exists")

type User struct {
	coll *mongo.Collection
}

func NewUser(client *mongo.Client) *User {
	return &User{coll: client.Database(configs.UserDBName).Collection(configs.UserCollName)}
}

// Get returns a user by given user id
func (u *User) Get(ctx context.Context, userId string) (*models.User, error) {
	var user models.User
	if err := u.coll.FindOne(ctx, bson.M{"_id": userId}).Decode(&user); err == mongo.ErrNoDocuments {
		return nil, ErrUserNotFound
	} else {
		return &user, nil
	}
}

// GetByUserName returns a user by given user name
func (u *User) GetByUserName(ctx context.Context, userName string) (*models.User, error) {
	var user models.User
	if err := u.coll.FindOne(ctx, bson.M{"user_name": userName}).Decode(&user); err == mongo.ErrNoDocuments {
		return nil, ErrUserNotFound
	} else {
		return &user, nil
	}
}

// GetAll returns a users by given user limit
func (u *User) GetAll(ctx context.Context, limit int64) (*[]models.User, error) {
	var users []models.User

	fr, err := u.coll.Find(ctx, bson.M{}, &options.FindOptions{Limit: &limit})
	if err = fr.All(ctx, &users); err != nil {
		return nil, err
	}

	return &users, nil
}

// Add creates a new user
func (u *User) Add(ctx context.Context, payload models.User) (string, error) {
	if _, err := u.coll.InsertOne(ctx, payload); err != nil {
		return "", err
	}

	return payload.Id, nil
}

// SetVerified updates user to verified by giver user id
func (u *User) SetVerified(ctx context.Context, userId string) error {
	ur, err := u.coll.UpdateByID(ctx, userId, bson.M{"is_verified": true})
	if err != nil {
		return err
	}
	if ur.MatchedCount == 0 {
		return ErrUserNotFound
	}
	return nil
}
