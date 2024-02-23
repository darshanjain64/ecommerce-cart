package controllers

import (
	"context"
	"ecommerce-cart/database"
	"ecommerce-cart/helper"
	"ecommerce-cart/models"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

var userCollection *mongo.Collection = database.OpenCollection(database.Client, "users")
var validate = validator.New()

func SignUp() gin.HandlerFunc{
	return func(c *gin.Context){
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)

		var user models.User

		if err:= c.BindJSON(&user); err!=nil{
			c.JSON(http.StatusBadRequest, gin.H{"error":err.Error()})
			return
		}

		validationErr:=validate.Struct(user)
		
		if validationErr!=nil{
			c.JSON(http.StatusBadRequest,gin.H{"error":validationErr.Error()})
			return
		}

		countEmail, err:=userCollection.CountDocuments(ctx, bson.M{"email":user.Email})
		defer cancel()

		if err!=nil{
			log.Panic(err)
			c.JSON(http.StatusInternalServerError, gin.H{"error":"error occured while checking mail"})
			return
		}

		password:=HashPassword(*user.Password)

		user.Password = &password

		countPhone, err:= userCollection.CountDocuments(ctx, bson.M{"phone":user.Phone})

		defer cancel()

		if err!=nil{
			log.Panic(err)
			c.JSON(http.StatusInternalServerError, gin.H{"error":"error occured while checking phone number"})
			return
		}
        
		if countEmail>0 || countPhone>0{
			c.JSON(http.StatusInternalServerError, gin.H{"error":"this email or phone number already exists"})
			return
		}

		user.Created_At, _ =time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		user.Updated_At, _ =time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		user.ID = primitive.NewObjectID()
		user.User_ID = user.ID.Hex()

		token, refreshToken,_:=helper.GenerateAllTokens(*user.Email, *user.First_Name, *user.Last_Name, *&user.User_ID)
		user.Token = &token
		user.RefreshToken = &refreshToken
		user.Address_Details = make([]models.Address,0)
		user.UserCart = make([]models.ProductUser,0)
		user.Order_Status = make([]models.Order, 0)

		_, insertErr:=userCollection.InsertOne(ctx, user)
		if insertErr!=nil{
			msg:=fmt.Sprintf("user was not created")
			c.JSON(http.StatusInternalServerError, gin.H{"error":msg})
			return
		}
        defer cancel()
		c.JSON(http.StatusCreated,"User created successfully")

	}
}

func Login() gin.HandlerFunc{
	return func(c *gin.Context){
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		var user models.User
		var foundUser models.User

		if err:= c.BindJSON(&user); err!=nil{
			c.JSON(http.StatusBadRequest, gin.H{"error":err.Error()})
			return
		}

		err:=userCollection.FindOne(ctx, bson.M{"email":user.Email}).Decode(&foundUser)
		defer cancel()
		if err!=nil{
			c.JSON(http.StatusInternalServerError, gin.H{"error":"user not found"})
			return
		}

		passwordIsValid, msg:= VerifyPassword(*user.Password, *foundUser.Password)
		defer cancel()
		if passwordIsValid!=true{
			c.JSON(http.StatusInternalServerError, gin.H{"error":msg})
			return
		}

		token, refreshToken, _:=helper.GenerateAllTokens(*foundUser.Email, *foundUser.First_Name, *foundUser.Last_Name, *&foundUser.User_ID)
		helper.UpdateAllTokens(token, refreshToken, foundUser.User_ID)

		c.JSON(http.StatusOK, foundUser)

	}
}

func HashPassword(password string) string{
    bytes, err:=bcrypt.GenerateFromPassword([]byte(password),14)
	if err!=nil{
		log.Panic(err)
	}
	return string(bytes)
}

func VerifyPassword(userPassword string, providePassword string) (bool,string){
	err:= bcrypt.CompareHashAndPassword([]byte(providePassword),[]byte(userPassword))
	check:=true
	msg:=""

	if err!=nil{
		msg=fmt.Sprintf("login or password is incorrect")
		check=false
	}
	return check, msg
}

