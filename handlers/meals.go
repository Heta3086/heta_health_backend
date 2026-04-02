package handlers

import (
	"net/http"
	"strings"
	"github.com/lib/pq"

	"heta_health_backend/config"

	"github.com/gin-gonic/gin"
)

func GetMeals(c *gin.Context) {

	diet := c.Query("diet")
	bmi := c.Query("bmi")
	allergies := c.Query("allergies") // comma separated

	var allergyList []string
	if allergies != "" {
		allergyList = strings.Split(allergies, ",")
	}

	// MAIN QUERY
	rows, err := config.DB.Query(`
	SELECT m.id, m.name, m.meal_type
	FROM meals m
	WHERE m.diet_type = $1
	AND m.id IN (
		SELECT meal_id FROM meal_bmi_categories WHERE bmi_category = $2
	)
	`, diet, bmi)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var result []gin.H

	for rows.Next() {
		var mealID int
		var name, mealType string

		rows.Scan(&mealID, &name, &mealType)

		// 🚫 Allergy Filter
		if len(allergyList) > 0 {

		for i := range allergyList {
			allergyList[i] = strings.ToLower(strings.TrimSpace(allergyList[i]))
		}

		var count int
		err := config.DB.QueryRow(`
			SELECT COUNT(*) 
			FROM meal_allergies ma
			JOIN allergies a ON ma.allergy_id = a.id
			WHERE ma.meal_id = $1 
			AND LOWER(a.name) = ANY($2)
		`, mealID, pq.Array(allergyList)).Scan(&count)

		if err == nil && count > 0 {
			continue // skip this meal
		}
	}

		// 🧪 Nutrition
		var calories int
		var protein, carbs, fats float64

		config.DB.QueryRow(`
			SELECT calories, protein, carbs, fats 
			FROM nutrition WHERE meal_id = $1
		`, mealID).Scan(&calories, &protein, &carbs, &fats)

		// 🧂 Ingredients
		ingRows, _ := config.DB.Query(`
			SELECT ingredient_name, quantity 
			FROM meal_ingredients WHERE meal_id = $1
		`, mealID)

		var ingredients []gin.H
		for ingRows.Next() {
			var name, qty string
			ingRows.Scan(&name, &qty)

			ingredients = append(ingredients, gin.H{
				"name":     name,
				"quantity": qty,
			})
		}

		// 📖 Recipes
		recRows, _ := config.DB.Query(`
			SELECT step_number, instruction 
			FROM recipes 
			WHERE meal_id = $1
			ORDER BY step_number
		`, mealID)

		var steps []gin.H
		for recRows.Next() {
			var step int
			var instruction string
			recRows.Scan(&step, &instruction)

			steps = append(steps, gin.H{
				"step":        step,
				"instruction": instruction,
			})
		}

		// 📦 Final Object
		result = append(result, gin.H{
			"id":         mealID,
			"name":       name,
			"meal_type":  mealType,
			"nutrition": gin.H{
				"calories": calories,
				"protein":  protein,
				"carbs":    carbs,
				"fats":     fats,
			},
			"ingredients": ingredients,
			"recipes":     steps,
		})
	}

	c.JSON(http.StatusOK, result)
}