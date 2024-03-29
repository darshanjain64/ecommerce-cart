package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type ProductUser struct {
	Product_ID       primitive.ObjectID   `bson:"_id"`
	Product_Name     *string              `json:"product_name" bson:"product_name"`
	Price            *uint64              `json:"price" bson:"price"`
	Rating           *uint8               `json:"rating" bson:"rating"`
	Image            *string              `json:"image" bson:"image"`
}