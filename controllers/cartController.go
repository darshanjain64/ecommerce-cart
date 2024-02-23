package controllers

import (
	"context"
	"ecommerce-cart/database"
	"ecommerce-cart/models"
	"errors"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type Application struct {
	prodCollection *mongo.Collection
	userCollection  *mongo.Collection
}

func NewApplication(prodCollection, userCollection *mongo.Collection) *Application{
	return &Application{
		prodCollection: prodCollection,
		userCollection: userCollection,
	}
}

func (app *Application)AddToCart() gin.HandlerFunc {
	return func(c *gin.Context) {
       productQueryId := c.Query("id")

	   if productQueryId == ""{
		log.Println("Product id is empty")

		_ = c.AbortWithError(http.StatusBadRequest, errors.New("product id is empty"))
		return
	   }

	   userQueryId := c.Query("userID")

	   if userQueryId == ""{
		log.Println("user id is empty")
		_ = c.AbortWithError(http.StatusBadRequest, errors.New("user id is empty"))
		return
	   }

	   productId, err := primitive.ObjectIDFromHex(productQueryId)

	   if err!=nil{
		log.Println(err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	   }

	   var ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)

	   defer cancel()

	   err = database.AddProductToCart(ctx, app.prodCollection, app.userCollection, productId, userQueryId)

	   if err!=nil{
		c.IndentedJSON(http.StatusInternalServerError, err)
		return
	   }
	   c.IndentedJSON(200, "Successfully added to the cart")
	}
}

func (app *Application)RemoveItem() gin.HandlerFunc {
	return func(c *gin.Context) {
       productQueryId := c.Query("id")

	   if productQueryId == ""{
		log.Println("Product id is empty")
		_ = c.AbortWithError(http.StatusBadRequest,errors.New("product id is empty"))
		return
	   }
 
	   userQueryId := c.Query("userID")

	   if userQueryId == ""{
		log.Println("user id is empty")
		_ = c.AbortWithError(http.StatusBadRequest, errors.New("user id is empty"))
		return
	   }

	   productId, err := primitive.ObjectIDFromHex(productQueryId)

	   if err!=nil{
		log.Println(err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	   }

	   var ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)

	   defer cancel()

	   err = database.RemoveCartItem(ctx, app.prodCollection, app.userCollection, productId, userQueryId)

	   if err!=nil{
		c.IndentedJSON(http.StatusInternalServerError,err)
		return
	   }
	   c.IndentedJSON(200,"Successfully removed from cart")
	}
}

func (app *Application)BuyFromCart() gin.HandlerFunc {
	return func(c *gin.Context) {
       userQueryId := c.Query("id")

	   if userQueryId ==""{
		log.Panic("user id is empty")
		_ = c.AbortWithError(http.StatusBadRequest, errors.New("user id is empty"))
		return
	   }

	   var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)

	   defer cancel()

	   err := database.BuyItemFromCart(ctx, app.userCollection, userQueryId)

	   if err!=nil{
		c.IndentedJSON(http.StatusInternalServerError, err)
	   }
	   c.IndentedJSON(200,"Successfully placed the order");
	}
}

func GetItemFromCart() gin.HandlerFunc {
	return func(c *gin.Context) {
         user_id := c.Query("id")

		 if user_id == ""{
			c.Header("Content-Type", "application/json")
			c.JSON(http.StatusNotFound, gin.H{"error":"invalid id"})
			c.Abort()
			return
		 }

		 usert_id, _ := primitive.ObjectIDFromHex(user_id)

		 var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		 defer cancel()

		 var filledCart models.User

		 err := userCollection.FindOne(ctx, bson.D{primitive.E{Key:"_id", Value: usert_id}}).Decode(&filledCart)

		 if err!=nil{
			log.Println(err)
			c.IndentedJSON(500, "not found")
			return
		 }

		 filter_match := bson.D{{Key:"$match", Value: bson.D{primitive.E{Key: "_id", Value: usert_id}}}}
		 unwind := bson.D{{Key:"$unwind", Value: bson.D{primitive.E{Key:"path", Value: "$user_cart"}}}}
		 grouping := bson.D{{Key: "$group",Value: bson.D{primitive.E{Key: "_id",Value: "$_id"},{Key:"total", Value: bson.D{primitive.E{Key: "$sum", Value: "$user_cart.price"}}}}}}
         
		 pointcursor, err := userCollection.Aggregate(ctx,mongo.Pipeline{filter_match,unwind,grouping})

		 if err!=nil{
			log.Println(err)
		 }

		 var listing []bson.M
		 if err = pointcursor.All(ctx, &listing); err!=nil{
			log.Println(err)
			c.AbortWithStatus(http.StatusInternalServerError)
		 }

        for _, json := range listing{
			c.IndentedJSON(200, json["total"])
			c.IndentedJSON(200,filledCart.UserCart)
		}

	}
}

func (app *Application)InstantBuy() gin.HandlerFunc {
	return func(c *gin.Context) {
        productQueryId := c.Query("id")

	   if productQueryId == ""{
		log.Println("Product id is empty")

		_ = c.AbortWithError(http.StatusBadRequest, errors.New("product id is empty"))
		return
	   }

	   userQueryId := c.Query("userID")

	   if userQueryId == ""{
		log.Println("user id is empty")
		_ = c.AbortWithError(http.StatusBadRequest, errors.New("user id is empty"))
		return
	   }

	   productId, err := primitive.ObjectIDFromHex(productQueryId)

	   if err!=nil{
		log.Println(err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	   }

	   var ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)

	   defer cancel()

	   err = database.InstantBuy(ctx, app.prodCollection, app.userCollection, productId, userQueryId)

	   if err!=nil{
		c.IndentedJSON(http.StatusInternalServerError, err)
		return
	   }
	   c.IndentedJSON(200, "Successfully placed the order")
	}
}

