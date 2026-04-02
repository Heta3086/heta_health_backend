package routes

import (
	"heta_health_backend/handlers"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(r *gin.Engine) {

	r.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "API Running"})
	})

	r.POST("/user", handlers.CreateUser)
	r.GET("/meals", handlers.GetMeals)
	r.GET("/recipes/:meal_id", handlers.GetRecipes)
}