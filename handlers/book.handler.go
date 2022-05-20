package handlers

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/agustadewa/book-system/models"
	"github.com/agustadewa/book-system/repo"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
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

func (h *BookHandler) RegisterEndpoints() {
	h.engine.POST("/book", h.addBook)
	h.engine.GET("/book/:book_id", h.getBook)
	h.engine.GET("/book/all", h.getAllBook)
	h.engine.DELETE("/book/:book_id", h.delete)
	h.engine.PUT("/book/updatestock/:book_id/:new_stock", h.updateBookStock)
	h.engine.POST("/book/update", h.updateBook)
}

func (h *BookHandler) addBook(c *gin.Context) {
	ctx := c.Request.Context()

	var addBookReq models.AddBookReq
	if err := c.BindJSON(&addBookReq); err != nil {
		var ve validator.ValidationErrors
		var allErrs string
		if errors.As(err, &ve) {
			var errStrs []string
			for _, fe := range ve {
				errStrs = append(errStrs, fmt.Sprintf("%s is %s", fe.Field(), fe.Tag()))
			}
			allErrs = strings.Join(errStrs, ", ")
		}
		c.JSON(http.StatusBadRequest, gin.H{"error": allErrs})
		return
	}

	addBook := models.AddBook{
		Price:       *addBookReq.Price,
		Qty:         *addBookReq.Qty,
		Name:        *addBookReq.Name,
		Author:      *addBookReq.Author,
		Publisher:   *addBookReq.Publisher,
		Category:    *addBookReq.Category,
		Language:    *addBookReq.Language,
		Description: *addBookReq.Description,
		Image:       *addBookReq.Image,
	}

	if addBook.Qty <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "quantity minimum is 1"})
		return
	}
	if addBook.Price <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "price minimum is 1"})
		return
	}

	// check existing book
	_, err := h.book.GetByName(ctx, addBook.Name)
	if err == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": repo.ErrBookExists.Error()})
		return
	}
	if err != nil && err != repo.ErrBookNotFound {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// add book
	addBookPayload := models.Book{
		Id:          primitive.NewObjectID().Hex(),
		Name:        addBook.Name,
		Author:      addBook.Author,
		Publisher:   addBook.Publisher,
		Category:    addBook.Category,
		Language:    addBook.Language,
		Price:       addBook.Price,
		Qty:         addBook.Qty,
		Description: addBook.Description,
		Image:       addBook.Image,
	}
	id, err := h.book.Add(ctx, addBookPayload)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "result": gin.H{"id": id}})
}

func (h *BookHandler) getBook(c *gin.Context) {
	ctx := c.Request.Context()

	bookId := c.Param("book_id")

	book, err := h.book.Get(ctx, bookId)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if book == nil {
		book = &models.Book{}
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

	if books == nil {
		bs := make([]models.Book, 0)
		books = &bs
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "result": books})
}

func (h *BookHandler) delete(c *gin.Context) {
	ctx := c.Request.Context()

	bookId := c.Param("book_id")

	if err := h.book.Delete(ctx, bookId); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "result": fmt.Sprintf("book %s has been deleted", bookId)})
}

func (h *BookHandler) updateBookStock(c *gin.Context) {
	ctx := c.Request.Context()

	bookId := c.Param("book_id")
	newStock := c.Param("new_stock")

	newStockInt, _ := strconv.ParseInt(newStock, 10, 64)

	if newStockInt < 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "number mustn't be negative"})
		return
	}

	if err := h.book.UpdateStock(ctx, bookId, newStockInt); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "result": fmt.Sprintf("book %s stock has been incremented by %v", bookId, newStockInt)})
}

func (h *BookHandler) updateBook(c *gin.Context) {
	ctx := c.Request.Context()

	var updateBook models.UpdateBookReq
	if err := c.BindJSON(&updateBook); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("error parsing request data: %s", err)})
		return
	}

	updatePayload := models.UpdateBook{
		Price:       updateBook.Price,
		Qty:         updateBook.Qty,
		Name:        updateBook.Name,
		Author:      updateBook.Author,
		Publisher:   updateBook.Publisher,
		Category:    updateBook.Category,
		Language:    updateBook.Language,
		Description: updateBook.Description,
		Image:       updateBook.Image,
	}

	if *updatePayload.Qty < 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Qty can't be lower than 0"})
		return
	}

	if *updatePayload.Price < 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Price can't be lower than 0"})
		return
	}

	if err := h.book.Update(ctx, *updateBook.Id, updatePayload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "result": fmt.Sprintf("Book %s has been updated", *updateBook.Id)})
}
