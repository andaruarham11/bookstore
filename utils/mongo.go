package utils

import (
	"context"
	"log"

	"github.com/agustadewa/book-system/configs"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func ConnectMongo(ctx context.Context) *mongo.Client {
	opt := options.Client().ApplyURI(configs.MongoUri)
	client, err := mongo.NewClient(opt)
	if err != nil {
		log.Fatalln("Can't create mongodb client: ", err)
	}

	if err = client.Connect(ctx); err != nil {
		log.Fatalln("Can't connect to mongodb server: ", err)
	}

	return client
}
