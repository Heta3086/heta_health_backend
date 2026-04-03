package routes

import (
	"heta_health_backend/handlers"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(r *gin.Engine) {

	r.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "API Running"})
	})

	// USER
	r.POST("/user", handlers.CreateUser)
	r.GET("/user/:id", handlers.GetUser)
	r.GET("/stats", handlers.GetStats)

	// MEALS
	r.GET("/meals", handlers.GetMeals)
	r.GET("/meals/:id", handlers.GetMealByID)

	// RECIPES
	r.GET("/recipes/:meal_id", handlers.GetRecipes)

	// ALLERGIES
	r.GET("/allergies", handlers.GetAllergies)
	r.GET("/diet-options", handlers.GetDietOptions)

	// FAVORITES
	r.POST("/favorites", handlers.AddFavorite)
	r.GET("/favorites/:user_id", handlers.GetFavorites)
	r.DELETE("/favorites", handlers.RemoveFavorite)

	// PLANNER
	r.POST("/planner", handlers.AddPlan)
	r.GET("/planner/:user_id", handlers.GetPlan)
	r.DELETE("/planner", handlers.RemovePlan)

	// AUTH
	r.POST("/signup", handlers.Signup)
	r.POST("/login", handlers.Login)
	r.POST("/logout", handlers.Logout)
}
