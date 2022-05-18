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

	handlers.NewUser(s, mClient).RegisterEndpoints()
	handlers.NewBook(s, mClient).RegisterEndpoints()
	handlers.NewOrder(s, mClient).RegisterEndpoints()
	handlers.NewPayment(s, mClient).RegisterEndpoints()

	utils.NewCronJob(mClient).DoCronJobTasks(ctx)

	if err := s.Run("0.0.0.0:4000"); err != nil {
		log.Fatalln("can't start server: ", err.Error())
	}
}
