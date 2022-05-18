package utils

import (
	"context"
	"log"
	"time"

	"github.com/agustadewa/book-system/models"
	"github.com/agustadewa/book-system/repo"
	"github.com/go-co-op/gocron"
	"go.mongodb.org/mongo-driver/mongo"
)

type cron struct {
	cron  *gocron.Scheduler
	order *repo.Order
	book  *repo.Book
}

func NewCronJob(mongoClient *mongo.Client) *cron {
	return &cron{
		cron:  gocron.NewScheduler(time.UTC),
		order: repo.NewOrder(mongoClient),
		book:  repo.NewBook(mongoClient),
	}
}

func (c *cron) removeAllIdleOrders(ctx context.Context) {
	orders, err := c.order.GetAllByStatus(ctx, models.WaitingForPayment, 100)
	if err != nil {
		log.Println("[CRON JOB ERROR] ", err.Error())
		return
	}

	for _, order := range *orders {
		orderTime, err := time.Parse(time.RFC3339, order.OrderTime)
		if err != nil {
			log.Println("[CRON JOB ERROR] ", err)
			return
		}

		gap := time.Now().Sub(orderTime)
		if gap > 30*time.Second {

			// update book quantity
			if err = c.book.UpdateStock(ctx, order.BookId, order.Qty); err != nil {
				log.Println("[CRON JOB ERROR] ", err)
				return
			}

			if err = c.order.Delete(ctx, order.Id); err != nil {
				log.Println("[CRON JOB ERROR] ", err.Error())
				return
			}

			log.Printf("[CRON JOB] order id %v is deleted\n", order.Id)
		}
	}
}

func (c *cron) DoCronJobTasks(ctx context.Context) {
	if _, err := c.cron.Every(5).Second().Do(func() {
		log.Println("[CRON JOB] doing cron job tasks")
		c.removeAllIdleOrders(ctx)
	}); err != nil {
		log.Println("[CRON JOB ERROR] ", err.Error())
		return
	}

	// c.cron.StartImmediately()
	c.cron.StartAsync()
}
