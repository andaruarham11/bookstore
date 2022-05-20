package models

type Book struct {
	Id          string `json:"id" bson:"_id"`
	Price       int64  `json:"price" bson:"price"`
	Qty         int64  `json:"qty" bson:"qty"`
	Name        string `json:"name" bson:"name"`
	Author      string `json:"author" bson:"author"`
	Publisher   string `json:"publisher" bson:"publisher"`
	Category    string `json:"category" bson:"category"`
	Language    string `json:"language" bson:"language"`
	Description string `json:"description" bson:"description" binding:"required"`
	Image       string `json:"image" bson:"image" binding:"required"`
}

type UpdateBookReq struct {
	Id          *string `json:"id" binding:"required"`
	Price       *int64  `json:"price"`
	Qty         *int64  `json:"qty"`
	Name        *string `json:"name"`
	Author      *string `json:"author"`
	Publisher   *string `json:"publisher"`
	Category    *string `json:"category"`
	Language    *string `json:"language"`
	Description *string `json:"description"`
	Image       *string `json:"image"`
}

type UpdateBook struct {
	Price       *int64  `bson:"price,omitempty"`
	Qty         *int64  `bson:"qty,omitempty"`
	Name        *string `bson:"name,omitempty"`
	Author      *string `bson:"author,omitempty"`
	Publisher   *string `bson:"publisher,omitempty"`
	Category    *string `bson:"category,omitempty"`
	Language    *string `bson:"language,omitempty"`
	Description *string `bson:"description,omitempty"`
	Image       *string `bson:"image,omitempty"`
}

type AddBook struct {
	Price       int64  `bson:"price"`
	Qty         int64  `bson:"qty"`
	Name        string `bson:"name"`
	Author      string `bson:"author"`
	Publisher   string `bson:"publisher"`
	Category    string `bson:"category"`
	Language    string `bson:"language"`
	Description string `bson:"description"`
	Image       string `bson:"image"`
}

type AddBookReq struct {
	Price       *int64  `json:"price" binding:"required"`
	Qty         *int64  `json:"qty" binding:"required"`
	Name        *string `json:"name" binding:"required"`
	Author      *string `json:"author" binding:"required"`
	Publisher   *string `json:"publisher" binding:"required"`
	Category    *string `json:"category" binding:"required"`
	Language    *string `json:"language" binding:"required"`
	Description *string `json:"description" binding:"required"`
	Image       *string `json:"image" binding:"required"`
}
