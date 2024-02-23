package controllers

import (
	"context"
	"ecommerce-cart/models"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func AddAddress() gin.HandlerFunc{
	return func(c *gin.Context){
       user_id := c.Query("id")
	   
	   if user_id==""{
		c.Header("Content-Type","application/json")
		c.JSON(http.StatusNotFound, gin.H{"error":"Invalid code"})
		c.Abort()
		return
	   }
	   address, err:= primitive.ObjectIDFromHex(user_id)
	   if err!=nil{
		c.IndentedJSON(500,"Internal Server Error")
	   }

	   var addresses models.Address

	   addresses.Address_ID = primitive.NewObjectID()

	   if err = c.BindJSON(&addresses); err!=nil{
		c.IndentedJSON(http.StatusNotAcceptable, err.Error())
	   }

	   var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
	   defer cancel()

       filter_match := bson.D{{Key:"$match",Value: bson.D{primitive.E{Key: "_id",Value: address}}}}
	   unwind := bson.D{{Key:"$unwind",Value: bson.D{primitive.E{Key: "path", Value:"$address" }}}}
	   group := bson.D{{Key:"$group",Value: bson.D{primitive.E{Key: "_id",Value: "address_id"}, {Key:"count",Value: bson.D{primitive.E{Key:"$sum", Value: 1}}}}}}

	   pointcursor, err := userCollection.Aggregate(ctx, mongo.Pipeline{filter_match, unwind, group})

	   if err!=nil{
		c.IndentedJSON(500, "Internal Server Error")
	   }
  
	   var addressinfo []bson.M

	   if err = pointcursor.All(ctx,&addressinfo);err!=nil{
		panic(err)
	   }
 
	   var size int32
	   for _, address_no := range addressinfo{
		count := address_no["count"]
		size = count.(int32)
	   }

	   if size<2{
		filter := bson.D{primitive.E{Key:"_id",Value: address}}
		update := bson.D{{Key: "$push", Value: bson.D{primitive.E{Key:"address", Value: addresses}}}}
	   _, err := userCollection.UpdateOne(ctx,filter,update)
	   if err!=nil{
		log.Println(err)
	   }
	   }else{
		c.IndentedJSON(400, "Not Allowed")
	   }
       defer cancel()
	   ctx.Done()
	}
}

func EditHomeAddress() gin.HandlerFunc{
	return func(c *gin.Context){
		user_id := c.Query("id")

		if user_id == ""{
			c.Header("Content-Type", "application/json")
			c.JSON(http.StatusNotFound, gin.H{"err":"invalid search index"})
			c.Abort()
			return
		}

		usert_id, err := primitive.ObjectIDFromHex(user_id)

		if err !=nil{
			c.IndentedJSON(500,"Internal Server Erorr")
			return
		}

		var editAddress models.Address
		if err = c.BindJSON(&editAddress);err!=nil{
			c.IndentedJSON(http.StatusBadRequest,err.Error())
		}

		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()
		filter := bson.D{primitive.E{Key:"_id",Value: usert_id}}
		update := bson.D{{Key: "$set", Value: bson.D{primitive.E{Key: "address.0.house_name", Value: editAddress.House}, {Key:"address.0.street_name", Value: editAddress.Street},{Key:"address.0.city_name", Value: editAddress.City},{Key:"address.0.pin_code", Value: editAddress.Pincode}}}}
		_, err = userCollection.UpdateOne(ctx, filter, update)
		if err!=nil{
			c.IndentedJSON(500,"Something went wrong")
			return
		}
		defer cancel()
		ctx.Done()
		c.IndentedJSON(200,"Successfully updated home address")
	}
}

func EditWorkAddress() gin.HandlerFunc{
	return func(c *gin.Context){
		user_id := c.Query("id")

		if user_id == ""{
			c.Header("Content-Type", "application/json")
			c.JSON(http.StatusNotFound, gin.H{"err":"invalid search index"})
			c.Abort()
			return
		}

		usert_id, err := primitive.ObjectIDFromHex(user_id)

		if err !=nil{
			c.IndentedJSON(500,"Internal Server Erorr")
			return
		}

		var editAddress models.Address
		if err = c.BindJSON(&editAddress);err!=nil{
			c.IndentedJSON(http.StatusBadRequest,err.Error())
		}

		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()
		filter := bson.D{primitive.E{Key:"_id",Value: usert_id}}
		update := bson.D{{Key: "$set", Value: bson.D{primitive.E{Key: "address.1.house_name", Value: editAddress.House}, {Key:"address.1.street_name", Value: editAddress.Street},{Key:"address.1.city_name", Value: editAddress.City},{Key:"address.1.pin_code", Value: editAddress.Pincode}}}}
		_, err = userCollection.UpdateOne(ctx, filter, update)
		if err!=nil{
			c.IndentedJSON(500,"Something went wrong")
			return
		}
		defer cancel()
		ctx.Done()
		c.IndentedJSON(200,"Successfully updated work address")
	}
}

func DeleteAddress() gin.HandlerFunc{
	return func(c *gin.Context){
		user_id := c.Query("id")

		if user_id == ""{
			c.Header("Content-Type", "application/json")
			c.JSON(http.StatusNotFound, gin.H{"err":"invalid search index"})
			c.Abort()
			return
		}
 
		addresses := make([]models.Address, 0)

		usert_id, err := primitive.ObjectIDFromHex(user_id)

		if err !=nil{
			c.IndentedJSON(500,"Internal Server Erorr")
			return
		}

		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()
		filter := bson.D{primitive.E{Key:"_id",Value: usert_id}}
		update := bson.D{{Key:"$set", Value:bson.D{primitive.E{Key:"address", Value:addresses}}}}
		_, err = userCollection.UpdateOne(ctx, filter, update)

		if err!=nil{
			c.IndentedJSON(404,"Wrong command")
			return
		}

		defer cancel()
		ctx.Done()
		c.IndentedJSON(200,"Successfully deleted")
	}
}