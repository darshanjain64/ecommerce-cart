package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Order struct {
	Order_ID primitive.ObjectID  `bson:"_id"`
	Order_Cart []ProductUser     `json: "order_list" bson:"list"` 
	Ordered_At time.Time         `json:"ordered_at" bson:"ordered_at"`
	Price   uint64               `json:"price" bson:"price"`
	Discount *int                `json:"discount" bson:"discount"`
	Payment_Method Payment       `json:"payment" bson:"payment"`
}