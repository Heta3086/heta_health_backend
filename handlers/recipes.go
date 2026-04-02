package handlers

import (
	"net/http"

	"heta_health_backend/config"

	"github.com/gin-gonic/gin"
)

func GetRecipes(c *gin.Context) {

	mealID := c.Param("meal_id")

	rows, _ := config.DB.Query(`
	SELECT step_number, instruction
	FROM recipes
	WHERE meal_id = $1
	ORDER BY step_number
	`, mealID)

	var steps []gin.H

	for rows.Next() {
		var step int
		var instruction string

		rows.Scan(&step, &instruction)

		steps = append(steps, gin.H{
			"step": step,
			"instruction": instruction,
		})
	}

	c.JSON(http.StatusOK, steps)
}