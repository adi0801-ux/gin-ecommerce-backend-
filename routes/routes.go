package routes

import (
	"ecommerce/controllers"
	"github.com/gin-gonic/gin"
)

func UserRoutes(incomingRoutes *gin.Engine) {
	incomingRoutes.POST("/users/signup", controllers.SignUp())
	incomingRoutes.POST("/users/login", controllers.Login())
	incomingRoutes.POST("/users/login", controllers.ProductViewerAdmin())
	incomingRoutes.POST("/users/login", controllers.SearchProduct())
	incomingRoutes.POST("/users/login", controllers.SearchProductByQuery())
}
