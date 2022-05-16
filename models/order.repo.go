package models

type OrderStatus string

const (
	Paid              OrderStatus = "PAID"
	Cancelled         OrderStatus = "CANCELLED"
	WaitingForPayment OrderStatus = "WAITING_FOR_PAYMENT"
	Declined          OrderStatus = "DECLINED"
	OnShipping        OrderStatus = "ON_SHIPPING"
)

func (o OrderStatus) IsPaid() bool {
	return o == Paid
}
func (o OrderStatus) IsCancelled() bool {
	return o == Cancelled
}
func (o OrderStatus) IsWaitingForPayment() bool {
	return o == WaitingForPayment
}
func (o OrderStatus) IsDeclined() bool {
	return o == Declined
}
func (o OrderStatus) IsOnShipping() bool {
	return o == OnShipping
}
func (o OrderStatus) String() string {
	return string(o)
}

type Order struct {
	Id        string      `json:"id" bson:"_id"`
	UserId    string      `json:"user_id" bson:"user_id"`
	BookId    string      `json:"book_id" bson:"book_id"`
	Qty       string      `json:"qty" bson:"qty"`
	OrderTime string      `json:"order_time" bson:"order_time"`
	Status    OrderStatus `json:"status" bson:"status"`
}

type AddOrder struct {
	UserId    string      `json:"user_id" bson:"user_id"`
	BookId    string      `json:"book_id" bson:"book_id"`
	Qty       string      `json:"qty" bson:"qty"`
	OrderTime string      `json:"order_time" bson:"order_time"`
	Status    OrderStatus `json:"status" bson:"status"`
}

type UpdateStatusOrder struct {
	OrderId string      `json:"order_id" bson:"order_id"`
	Status  OrderStatus `json:"status" bson:"status"`
}
