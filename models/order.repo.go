package models

import "errors"

type OrderStatus string

var ErrUnknownOrderStatus = errors.New("unknown order status")

const (
	Paid              OrderStatus = "PAID"
	Cancelled         OrderStatus = "CANCELLED"
	WaitingForPayment OrderStatus = "WAITING_FOR_PAYMENT"
	Declined          OrderStatus = "DECLINED"
	OnShipping        OrderStatus = "ON_SHIPPING"
)

func IsValidOrderStatus(status string) (OrderStatus, error) {
	switch status {
	case Paid.String():
		break
	case Cancelled.String():
		break
	case WaitingForPayment.String():
		break
	case Declined.String():
		break
	case OnShipping.String():
		break
	default:
		return "", ErrUnknownOrderStatus
	}

	return OrderStatus(status), nil
}

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
	Id         string      `json:"id" bson:"_id"`
	UserId     string      `json:"user_id" bson:"user_id"`
	BookId     string      `json:"book_id" bson:"book_id"`
	OrderTime  string      `json:"order_time" bson:"order_time"`
	Status     OrderStatus `json:"status" bson:"status"`
	Qty        int64       `json:"qty" bson:"qty"`
	TotalPrice int64       `json:"total_price" bson:"total_price"`
}

type AddOrder struct {
	UserId string `json:"user_id" bson:"user_id" binding:"required"`
	BookId string `json:"book_id" bson:"book_id" binding:"required"`
	Qty    int64  `json:"qty" bson:"qty" binding:"required"`
}

type UpdateStatusOrder struct {
	OrderId string      `json:"order_id" bson:"order_id"`
	Status  OrderStatus `json:"status" bson:"status"`
}
