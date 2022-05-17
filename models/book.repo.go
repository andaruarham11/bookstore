package models

type Book struct {
	Price     int64  `json:"price" bson:"price"`
	Qty       int64  `json:"qty" bson:"qty"`
	Id        string `json:"id" bson:"_id"`
	Name      string `json:"name" bson:"name"`
	Author    string `json:"author" bson:"author"`
	Publisher string `json:"publisher" bson:"publisher"`
	Category  string `json:"category" bson:"category"`
	Language  string `json:"language" bson:"language"`
	Image     string `json:"image" bson:"image"`
}

type AddBook struct {
	Price     int64  `json:"price" bson:"price" binding:"required"`
	Qty       int64  `json:"qty" bson:"qty" binding:"required"`
	Name      string `json:"name" bson:"name" binding:"required"`
	Author    string `json:"author" bson:"author" binding:"required"`
	Publisher string `json:"publisher" bson:"publisher" binding:"required"`
	Category  string `json:"category" bson:"category" binding:"required"`
	Language  string `json:"language" bson:"language" binding:"required"`
	Image     string `json:"image" bson:"image" binding:"required"`
}
