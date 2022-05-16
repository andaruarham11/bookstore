package models

type Book struct {
	Id        string `json:"id" bson:"_id"`
	Name      string `json:"name" bson:"name"`
	Author    string `json:"author" bson:"author"`
	Publisher string `json:"publisher" bson:"publisher"`
	Category  string `json:"category" bson:"category"`
	Language  string `json:"language" bson:"language"`
	Price     int64  `json:"price" bson:"price"`
	Qty       int64  `json:"qty" bson:"qty"`
}

type AddBook struct {
	Name      string `json:"name" bson:"name"`
	Author    string `json:"author" bson:"author"`
	Publisher string `json:"publisher" bson:"publisher"`
	Category  string `json:"category" bson:"category"`
	Language  string `json:"language" bson:"language"`
	Price     int64  `json:"price" bson:"price"`
	Qty       int64  `json:"qty" bson:"qty"`
}
