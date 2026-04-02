package handlers

import (
	"net/http"

	"heta_health_backend/config"
	"heta_health_backend/utils"

	"github.com/gin-gonic/gin"
)

func CreateUser(c *gin.Context) {

	var input struct {
		Name   string
		Age    int
		Gender string
		Height float64
		Weight float64
		Diet   string
	}

	c.BindJSON(&input)

	bmi := utils.CalculateBMI(input.Weight, input.Height)
	category := utils.GetBMICategory(bmi)

	_, err := config.DB.Exec(
		"INSERT INTO users (name, age, gender, height_cm, weight_kg, diet_preference) VALUES ($1,$2,$3,$4,$5,$6)",
		input.Name, input.Age, input.Gender, input.Height, input.Weight, input.Diet,
	)

	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":  "User created",
		"bmi":      bmi,
		"category": category,
	})
}