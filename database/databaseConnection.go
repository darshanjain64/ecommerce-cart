package database

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func DBinstance() *mongo.Client{
	serverAPI := options.ServerAPI(options.ServerAPIVersion1)
	opts := options.Client().ApplyURI("mongodb+srv://darshanr94dj:Hello%402024@cluster0.o9oelrd.mongodb.net/?retryWrites=true&w=majority&appName=Cluster0").SetServerAPIOptions(serverAPI)
  
	client, err := mongo.Connect(context.TODO(), opts)
  if err != nil {
    panic(err)
  }
  
  // Send a ping to confirm a successful connection
  if err := client.Database("admin").RunCommand(context.TODO(), bson.D{{Key: "ping", Value: 1}}).Err(); err != nil {
    panic(err)
  }
  fmt.Println("Pinged your deployment. You successfully connected to MongoDB!")
	return client
}


  var Client *mongo.Client = DBinstance()

  func OpenCollection(client *mongo.Client, collectionName string) *mongo.Collection{
    var collection *mongo.Collection = client.Database("ecommerce").Collection(collectionName)

	return collection
  }
