package middleware

import (
	"ecommerce-cart/helper"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

func Authentication() gin.HandlerFunc{
	return func(c *gin.Context){
        clientToken:=c.Request.Header.Get("token")

		if clientToken==""{
			c.JSON(http.StatusInternalServerError,gin.H{"error":fmt.Sprintf("Not Authorization Header Provided")})
			c.Abort()
			return
		}
		claims,err:=helper.ValidateToken(clientToken)
		if err!=""{
			c.JSON(http.StatusInternalServerError,gin.H{"error":err})
			c.Abort()
			return
		}

		c.Set("email",claims.Email)
		c.Set("uid",claims.Uid)
		c.Next()
	}
}