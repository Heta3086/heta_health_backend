package handlers

import (
	"net/http"

	"heta_health_backend/config"
	"heta_health_backend/utils"

	"github.com/gin-gonic/gin"
)

func CreateUser(c *gin.Context) {

	var input struct {
		Name   string  `json:"name"`
		Age    int     `json:"age"`
		Gender string  `json:"gender"`
		Height float64 `json:"height"`
		Weight float64 `json:"weight"`
		Diet   string  `json:"diet"`
	}

	// ✅ Bind JSON
	if err := c.BindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	// ✅ VALIDATION (IMPORTANT 🔥)
	if input.Age < 18 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "BMI only valid for age 18+",
		})
		return
	}

	if input.Height <= 0 || input.Weight <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid height or weight",
		})
		return
	}

	// ✅ BMI calculation
	bmi := utils.CalculateBMI(input.Weight, input.Height)
	category := utils.GetBMICategory(bmi)

	// ✅ INSERT + RETURN ID
	var userID int
	err := config.DB.QueryRow(
		"INSERT INTO users (name, age, gender, height_cm, weight_kg, diet_preference) VALUES ($1,$2,$3,$4,$5,$6) RETURNING id",
		input.Name, input.Age, input.Gender, input.Height, input.Weight, input.Diet,
	).Scan(&userID)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// ✅ RESPONSE
	c.JSON(http.StatusOK, gin.H{
		"user_id":  userID,
		"bmi":      bmi,
		"category": category,
	})
}

func GetStats(c *gin.Context) {
	var stats struct {
		RecipeCount int `json:"recipeCount"`
		UserCount   int `json:"userCount"`
		MealCount   int `json:"mealCount"`
	}

	err := config.DB.QueryRow(`
		SELECT
			(SELECT COUNT(DISTINCT meal_id) FROM recipes),
			(SELECT COUNT(*) FROM users),
			(SELECT COUNT(*) FROM meals)
	`).Scan(&stats.RecipeCount, &stats.UserCount, &stats.MealCount)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, stats)
}

func GetUser(c *gin.Context) {
	id := c.Param("id")

	var uid int
	var name, gender, diet string
	var age int
	var height, weight float64

	err := config.DB.QueryRow(
		"SELECT id, name, age, gender, height_cm, weight_kg, diet_preference FROM users WHERE id=$1",
		id,
	).Scan(&uid, &name, &age, &gender, &height, &weight, &diet)

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	user := gin.H{
		"id":     uid,
		"name":   name,
		"age":    age,
		"gender": gender,
		"height": height,
		"weight": weight,
		"diet":   diet,
	}

	c.JSON(http.StatusOK, user)
}
