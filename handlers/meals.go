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
	allergies := c.Query("allergies")

	var allergyList []string
	if allergies != "" {
		allergyList = strings.Split(allergies, ",")
		for i := range allergyList {
			allergyList[i] = strings.ToLower(strings.TrimSpace(allergyList[i]))
		}
	}

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
	defer rows.Close()

	var result []gin.H

	for rows.Next() {
		var mealID int
		var name, mealType string

		rows.Scan(&mealID, &name, &mealType)

		// 🚫 Allergy Filter
		if len(allergyList) > 0 {
			var count int

			err := config.DB.QueryRow(`
				SELECT COUNT(*) 
				FROM meal_allergies ma
				JOIN allergies a ON ma.allergy_id = a.id
				WHERE ma.meal_id = $1 
				AND LOWER(a.name) = ANY($2)
			`, mealID, pq.Array(allergyList)).Scan(&count)

			if err == nil && count > 0 {
				continue
			}
		}

		// 🧪 Nutrition
		var calories int
		var protein, carbs, fats float64

		err := config.DB.QueryRow(`
			SELECT calories, protein, carbs, fats 
			FROM nutrition WHERE meal_id = $1
		`, mealID).Scan(&calories, &protein, &carbs, &fats)

		if err != nil {
			continue
		}

		// 🧂 Ingredients
		ingRows, err := config.DB.Query(`
			SELECT ingredient_name, quantity 
			FROM meal_ingredients WHERE meal_id = $1
		`, mealID)

		if err != nil {
			continue
		}
		defer ingRows.Close()

		var ingredients []gin.H
		for ingRows.Next() {
			var ingName, qty string
			ingRows.Scan(&ingName, &qty)

			ingredients = append(ingredients, gin.H{
				"name":     ingName,
				"quantity": qty,
			})
		}

		// 📖 Recipes
		recRows, err := config.DB.Query(`
			SELECT step_number, instruction 
			FROM recipes 
			WHERE meal_id = $1
			ORDER BY step_number
		`, mealID)

		if err != nil {
			continue
		}
		defer recRows.Close()

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

		// 📦 Final Response
		result = append(result, gin.H{
			"id":        mealID,
			"name":      name,
			"meal_type": mealType,
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

func GetMealByID(c *gin.Context) {
	id := c.Param("id")

	var meal struct {
		ID       int    `json:"id"`
		Name     string `json:"name"`
		MealType string `json:"meal_type"`
	}

	err := config.DB.QueryRow(
		"SELECT id, name, meal_type FROM meals WHERE id=$1", id,
	).Scan(&meal.ID, &meal.Name, &meal.MealType)

	if err != nil {
		c.JSON(404, gin.H{"error": "Meal not found"})
		return
	}

	// ✅ FIXED Nutrition (IMPORTANT 🔥)
	var calories int
	var protein, carbs, fats float64

	err = config.DB.QueryRow(
		"SELECT calories, protein, carbs, fats FROM nutrition WHERE meal_id=$1",
		id,
	).Scan(&calories, &protein, &carbs, &fats)

	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to fetch nutrition"})
		return
	}

	nutrition := gin.H{
		"calories": calories,
		"protein":  protein,
		"carbs":    carbs,
		"fats":     fats,
	}

	// 🧂 Ingredients
	rows, _ := config.DB.Query(
		"SELECT ingredient_name, quantity FROM meal_ingredients WHERE meal_id=$1",
		id,
	)
	defer rows.Close()

	var ingredients []gin.H
	for rows.Next() {
		var name, qty string
		rows.Scan(&name, &qty)

		ingredients = append(ingredients, gin.H{
			"name":     name,
			"quantity": qty,
		})
	}

	// 📖 Recipes
	rows2, _ := config.DB.Query(
		"SELECT step_number, instruction FROM recipes WHERE meal_id=$1 ORDER BY step_number",
		id,
	)
	defer rows2.Close()

	var recipes []gin.H
	for rows2.Next() {
		var step int
		var instruction string
		rows2.Scan(&step, &instruction)

		recipes = append(recipes, gin.H{
			"step":        step,
			"instruction": instruction,
		})
	}

	c.JSON(200, gin.H{
		"meal":        meal,
		"nutrition":   nutrition,
		"ingredients": ingredients,
		"recipes":     recipes,
	})
}