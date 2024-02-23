package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Address struct {
	Address_ID primitive.ObjectID  `bson:"_id"`
	House      *string             `json:"house" bson:"house"`
	Street     *string             `json:"steet" bson:"street"`
	City       *string             `json:"city" bson:"city"`
	Pincode    *string             `json:"pincode" bson:"pincode"`
}