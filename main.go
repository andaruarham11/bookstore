package main

import (
	"context"
	"log"

	"github.com/agustadewa/book-system/handlers"
	"github.com/agustadewa/book-system/utils"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	ctx := context.Background()

	s := gin.Default()
	s.Use(cors.New(cors.Config{AllowOrigins: []string{"*"}, AllowCredentials: true}))

	mClient := utils.ConnectMongo(ctx)

	handlers.NewUser(s, mClient).RegisterEndpoint()
	handlers.NewBook(s, mClient).RegisterEndpoint()
	handlers.NewOrder(s, mClient).RegisterEndpoint()
	handlers.NewPayment(s, mClient).RegisterEndpoint()

	if err := s.Run("0.0.0.0:4000"); err != nil {
		log.Fatalln("can't start server: ", err.Error())
	}
}
