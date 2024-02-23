package routes

import (
	"ecommerce-cart/controllers"
	"github.com/gin-gonic/gin"
)

func AddressRoutes(incomingRoutes *gin.Engine){
	incomingRoutes.POST("address/add", controllers.AddAddress())
	incomingRoutes.PUT("address/edithome", controllers.EditHomeAddress())
	incomingRoutes.PUT("address/editwork", controllers.EditWorkAddress())
	incomingRoutes.GET("address/delete", controllers.DeleteAddress())
}