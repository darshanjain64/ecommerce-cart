package routes

import (
	"ecommerce-cart/controllers"
	"github.com/gin-gonic/gin"
)

func UserRoutes(incomingRoutes *gin.Engine){
	incomingRoutes.POST("users/signup",controllers.SignUp())
	incomingRoutes.POST("users/login",controllers.Login())
}