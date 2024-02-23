package routes

import (
	controllers "ecommerce-cart/controllers"
	"ecommerce-cart/database"

	"github.com/gin-gonic/gin"
)

var app = controllers.NewApplication(database.OpenCollection(database.Client,"products"),database.OpenCollection(database.Client,"users")) 

func CartRoutes(incomingRoutes *gin.Engine){
	incomingRoutes.GET("cart/addtocart",app.AddToCart())
	incomingRoutes.GET("cart/removeitem",app.RemoveItem())
	incomingRoutes.GET("cart/cartcheckout",app.BuyFromCart())
	incomingRoutes.GET("cart/instantbuy",app.InstantBuy())
}