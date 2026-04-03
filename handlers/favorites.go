package handlers

import (
	"heta_health_backend/config"
	"github.com/gin-gonic/gin"
)

func AddFavorite(c *gin.Context) {
	var data struct {
		UserID int `json:"user_id"`
		MealID int `json:"meal_id"`
	}

	c.BindJSON(&data)

	config.DB.Exec(
		"INSERT INTO favorites (user_id, meal_id) VALUES ($1,$2)",
		data.UserID, data.MealID,
	)

	c.JSON(200, gin.H{"message": "Added"})
}

func GetFavorites(c *gin.Context) {
	userID := c.Param("user_id")

	rows, _ := config.DB.Query(
		"SELECT m.id, m.name FROM favorites f JOIN meals m ON f.meal_id=m.id WHERE f.user_id=$1",
		userID,
	)

	var meals []map[string]interface{}
	for rows.Next() {
		var id int
		var name string
		rows.Scan(&id, &name)

		meals = append(meals, gin.H{
			"id": id,
			"name": name,
		})
	}

	c.JSON(200, meals)
}

func RemoveFavorite(c *gin.Context) {
	var data struct {
		UserID int `json:"user_id"`
		MealID int `json:"meal_id"`
	}

	c.BindJSON(&data)

	config.DB.Exec(
		"DELETE FROM favorites WHERE user_id=$1 AND meal_id=$2",
		data.UserID, data.MealID,
	)

	c.JSON(200, gin.H{"message": "Removed"})
}