package handlers

import (
	"database/sql"
	"net/http"

	"heta_health_backend/config"
	"heta_health_backend/utils"

	"github.com/gin-gonic/gin"
)

func CreateUser(c *gin.Context) {

	var input struct {
		AuthUserID int     `json:"auth_user_id"`
		Name       string  `json:"name"`
		Age        int     `json:"age"`
		Gender     string  `json:"gender"`
		Height     float64 `json:"height"`
		Weight     float64 `json:"weight"`
		Diet       string  `json:"diet"`
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

	if input.AuthUserID <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "auth_user_id is required"})
		return
	}

	// Keep a single profile row per auth account.
	var userID int
	err := config.DB.QueryRow(
		`UPDATE users
		 SET name=$1, age=$2, gender=$3, height_cm=$4, weight_kg=$5, diet_preference=$6
		 WHERE auth_user_id=$7
		 RETURNING id`,
		input.Name, input.Age, input.Gender, input.Height, input.Weight, input.Diet, input.AuthUserID,
	).Scan(&userID)

	if err == sql.ErrNoRows {
		err = config.DB.QueryRow(
			"INSERT INTO users (auth_user_id, name, age, gender, height_cm, weight_kg, diet_preference) VALUES ($1,$2,$3,$4,$5,$6,$7) RETURNING id",
			input.AuthUserID, input.Name, input.Age, input.Gender, input.Height, input.Weight, input.Diet,
		).Scan(&userID)
	}

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
			(SELECT COUNT(*) FROM auth_users),
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
