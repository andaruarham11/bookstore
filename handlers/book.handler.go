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

func NewBook(engine *gin.Engine, client *mongo.Client) *BookHandler {
	return &BookHandler{
		engine: engine,
		book:   repo.NewBook(client),
	}
}

type BookHandler struct {
	engine *gin.Engine
	book   *repo.Book
}

func (h *BookHandler) RegisterEndpoint() {
	h.engine.POST("/book", h.addBook)
	h.engine.GET("/book/:book_id", h.getBook)
	h.engine.GET("/book/all", h.getAllBook)
}

func (h *BookHandler) addBook(c *gin.Context) {
	ctx := c.Request.Context()

	var addBook models.AddBook
	if err := c.BindJSON(&addBook); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("error parsing request data: %s", err)})
		return
	}

	// get book
	_, err := h.book.GetByName(ctx, addBook.Name)
	if err == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": repo.ErrBookExists.Error()})
		return
	}
	if err != nil && err != repo.ErrBookNotFound {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}

	// add book
	addBookPayload := models.Book{
		Id:        primitive.NewObjectID().Hex(),
		Name:      addBook.Name,
		Author:    addBook.Author,
		Publisher: addBook.Publisher,
		Category:  addBook.Category,
		Language:  addBook.Language,
		Price:     addBook.Price,
		Qty:       addBook.Qty,
		Image:     addBook.Image,
	}
	id, err := h.book.Add(ctx, addBookPayload)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"result":  gin.H{"id": id},
	})
}

func (h *BookHandler) getBook(c *gin.Context) {
	ctx := c.Request.Context()

	bookId := c.Param("book_id")

	book, err := h.book.Get(ctx, bookId)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "result": book})
}

func (h *BookHandler) getAllBook(c *gin.Context) {
	ctx := c.Request.Context()

	limitStr := c.Request.URL.Query().Get("limit")
	limit, _ := strconv.ParseInt(limitStr, 10, 64)

	if limit < 10 || limit > 100 {
		limit = 10
	}
	books, err := h.book.GetAll(ctx, limit)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "result": books})
}
