package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	ID              primitive.ObjectID   `json:"_id" bson:"_id"`
	First_Name      *string              `json:"first_name" validate:"required"`
	Last_Name       *string              `json:"last_name" validate:"required"`
	Password        *string              `json:"password" validate:"required"`
	Email           *string              `json:"email" validate:"required"`
	Phone           *string              `json:"phone"`
	Token           *string              `json:"token"`
	RefreshToken    *string              `json:"refresh_token"`
	Created_At      time.Time            `json:"created_at"`
	Updated_At      time.Time            `json:"updated_at"`
	User_ID         string               `json:"user_id"`
	UserCart        []ProductUser        `json:"user_cart" bson:"user_cart"`
	Address_Details []Address            `json:"address" bson:"address"`
	Order_Status    []Order              `json:"order" bson:"orders"` 
}