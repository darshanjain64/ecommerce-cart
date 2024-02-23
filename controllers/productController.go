package controllers

import (
	"context"
	"ecommerce-cart/database"
	"ecommerce-cart/models"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var productCollection *mongo.Collection = database.OpenCollection(database.Client, "products")

func AddProducts() gin.HandlerFunc {
	return func(c *gin.Context) {
       var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
	   var products models.Product
	   defer cancel()

	   if err:=c.BindJSON(&products);err!=nil{
		c.JSON(http.StatusBadRequest, gin.H{"error":err.Error()})
		return
	   }

	   products.Product_ID = primitive.NewObjectID()
	   _, posterr := productCollection.InsertOne(ctx, products)
	   if posterr!=nil{
		c.JSON(http.StatusInternalServerError, gin.H{"err":"not inserted"})
		return
	   }
	   defer cancel()
	   c.JSON(http.StatusOK,"Product added")
	}
}

func SearchProduct() gin.HandlerFunc {
	return func(c *gin.Context) {
            var productList []models.Product
			var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
			defer cancel()
 
			 cursor, err := productCollection.Find(ctx, bson.D{{}})

			 if err!=nil{
				c.IndentedJSON(http.StatusInternalServerError, "something went wrong, please try again later")
				return
			 }

			 err = cursor.All(ctx, &productList)
			 if err!=nil{
				log.Println(err)
				c.AbortWithStatus(http.StatusInternalServerError)
				return
			 }
			 defer cursor.Close(ctx)

			 if err := cursor.Err();err!=nil{
				log.Println(err)
				c.IndentedJSON(400, "invalid")
				return
			 }
			 defer cancel()

			 c.IndentedJSON(200,productList)

	}
}

func SearchProductByQuery() gin.HandlerFunc {
	return func(c *gin.Context) {
        var searchProducts []models.Product

		queryParam := c.Query("name")

		if queryParam ==""{
			log.Println("query is empty")
			c.Header("Content-Type","application/json")
			c.JSON(http.StatusNotFound, gin.H{"Error":"Invalid search index"})
			c.Abort()
			return
		}

		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		searchquerydb, err := productCollection.Find(ctx, bson.M{"product_name": bson.M{"$regex":queryParam}})

		if err!=nil{
			c.IndentedJSON(404, "something went wrong while fetching the data")
			return
		}

		err = searchquerydb.All(ctx, &searchProducts)
		if err!=nil{
			log.Println(err)
			c.IndentedJSON(400, "invalid")
			return
		}

		defer searchquerydb.Close(ctx)

		if err := searchquerydb.Err(); err!=nil{
			log.Println(err)
			c.IndentedJSON(400, "invalid request")
			return
		}

		defer cancel()
		c.IndentedJSON(200, searchProducts)
	}
}