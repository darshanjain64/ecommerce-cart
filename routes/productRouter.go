package routes


import (
	"ecommerce-cart/controllers"
	"github.com/gin-gonic/gin"
)

func ProductRoutes(incomingRoutes *gin.Engine){
	incomingRoutes.POST("admin/addProducts",controllers.AddProducts())
	incomingRoutes.GET("products/productView",controllers.SearchProduct())
	incomingRoutes.GET("products/search",controllers.SearchProductByQuery())
}