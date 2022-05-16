package models

type Payment struct {
	Id      string `json:"id" bson:"_id"`
	UserId  string `json:"user_id" bson:"user_id"`
	OrderId string `json:"order_id" bson:"order_id"`
	Receipt string `json:"receipt" bson:"receipt"`
}

type AddPayment struct {
	UserId  string `json:"user_id" bson:"user_id"`
	OrderId string `json:"order_id" bson:"order_id"`
	Receipt string `json:"receipt" bson:"receipt"`
}
