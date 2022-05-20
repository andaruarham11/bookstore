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

var ErrBookNotFound = errors.New("book not found")
var ErrBookExists = errors.New("book already exists")

type Book struct {
	coll *mongo.Collection
}

func NewBook(client *mongo.Client) *Book {
	return &Book{coll: client.Database(configs.BookDBName).Collection(configs.BookCollName)}
}

// Get returns a book by given book id
func (b *Book) Get(ctx context.Context, bookId string) (*models.Book, error) {
	var book models.Book
	if err := b.coll.FindOne(ctx, bson.M{"_id": bookId}).Decode(&book); err == mongo.ErrNoDocuments {
		return nil, ErrBookNotFound
	} else {
		return &book, nil
	}
}

// GetByName returns a book by given book name
func (b *Book) GetByName(ctx context.Context, bookName string) (*models.Book, error) {
	var book models.Book
	if err := b.coll.FindOne(ctx, bson.M{"name": bookName}).Decode(&book); err == mongo.ErrNoDocuments {
		return nil, ErrBookNotFound
	} else {
		return &book, nil
	}
}

// GetAll returns a books by given book limit
func (b *Book) GetAll(ctx context.Context, limit int64) (*[]models.Book, error) {
	var books []models.Book

	// fr, err := b.coll.Find(ctx, bson.M{"qty": bson.M{"$ne": 0}}, &options.FindOptions{
	// 	Limit: &limit,
	// })
	fr, err := b.coll.Find(ctx, bson.M{}, &options.FindOptions{
		Limit: &limit,
	})
	if err = fr.All(ctx, &books); err != nil {
		return nil, err
	}

	return &books, nil
}

// Add creates a new book
func (b *Book) Add(ctx context.Context, payload models.Book) (string, error) {
	if _, err := b.coll.InsertOne(ctx, payload); err != nil {
		return "", err
	}

	return payload.Id, nil
}

// Delete deletes a book
func (b *Book) Delete(ctx context.Context, bookId string) error {
	dr, err := b.coll.DeleteOne(ctx, bson.M{"_id": bookId})
	if err != nil {
		return err
	}
	if dr.DeletedCount == 0 {
		return ErrBookNotFound
	}
	return nil
}

// UpdateStock updates a book stock
func (b *Book) UpdateStock(ctx context.Context, bookId string, newStock int64) error {
	ur, err := b.coll.UpdateByID(ctx, bookId, bson.M{"$inc": bson.M{"qty": newStock}})
	if err != nil {
		return err
	}
	if ur.MatchedCount == 0 {
		return ErrBookNotFound
	}
	return nil
}

// Update updates a book
func (b *Book) Update(ctx context.Context, bookId string, updatePayload models.UpdateBook) error {
	ur, err := b.coll.UpdateByID(ctx, bookId, bson.M{"$set": updatePayload})
	if err != nil {
		return err
	}
	if ur.MatchedCount == 0 {
		return ErrBookNotFound
	}
	return nil
}
