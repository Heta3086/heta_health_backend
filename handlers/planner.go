package handlers

import (
	"heta_health_backend/config"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

var validWeekdays = map[string]string{
	"monday":    "Monday",
	"tuesday":   "Tuesday",
	"wednesday": "Wednesday",
	"thursday":  "Thursday",
	"friday":    "Friday",
	"saturday":  "Saturday",
	"sunday":    "Sunday",
}

func normalizeWeekday(day string) (string, bool) {
	normalized, ok := validWeekdays[strings.ToLower(strings.TrimSpace(day))]
	return normalized, ok
}

func AddPlan(c *gin.Context) {
	var data struct {
		UserID int    `json:"user_id"`
		MealID int    `json:"meal_id"`
		Day    string `json:"day"`
	}

	if err := c.BindJSON(&data); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid planner payload"})
		return
	}

	normalizedDay, validDay := normalizeWeekday(data.Day)

	if data.UserID <= 0 || data.MealID <= 0 || !validDay {
		c.JSON(http.StatusBadRequest, gin.H{"error": "user_id, meal_id and day are required"})
		return
	}

	_, err := config.DB.Exec(
		"DELETE FROM planner WHERE user_id=$1 AND day=$2",
		data.UserID, normalizedDay,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update existing plan"})
		return
	}

	_, err = config.DB.Exec(
		"INSERT INTO planner (user_id, meal_id, day) VALUES ($1,$2,$3)",
		data.UserID, data.MealID, normalizedDay,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add plan"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Added"})
}

func GetPlan(c *gin.Context) {
	userID := c.Param("user_id")

	rows, err := config.DB.Query(
		`SELECT p.day, m.id, m.name, m.meal_type, COALESCE(n.calories, 0)
		 FROM planner p
		 JOIN meals m ON p.meal_id = m.id
		 LEFT JOIN nutrition n ON n.meal_id = m.id
		 WHERE p.user_id = $1
		 ORDER BY CASE p.day
		 	WHEN 'Monday' THEN 1
		 	WHEN 'Tuesday' THEN 2
		 	WHEN 'Wednesday' THEN 3
		 	WHEN 'Thursday' THEN 4
		 	WHEN 'Friday' THEN 5
		 	WHEN 'Saturday' THEN 6
		 	WHEN 'Sunday' THEN 7
		 	ELSE 8
		 END`,
		userID,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch planner"})
		return
	}
	defer rows.Close()

	var plans []map[string]interface{}
	for rows.Next() {
		var day, name, mealType string
		var mealID, calories int
		err := rows.Scan(&day, &mealID, &name, &mealType, &calories)
		if err != nil {
			continue
		}

		plans = append(plans, gin.H{
			"day": day,
			"meal": gin.H{
				"id":       mealID,
				"name":     name,
				"type":     mealType,
				"calories": calories,
			},
		})
	}

	c.JSON(http.StatusOK, plans)
}

func RemovePlan(c *gin.Context) {
	var data struct {
		UserID int    `json:"user_id"`
		Day    string `json:"day"`
	}

	if err := c.BindJSON(&data); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid planner payload"})
		return
	}

	normalizedDay, validDay := normalizeWeekday(data.Day)

	if data.UserID <= 0 || !validDay {
		c.JSON(http.StatusBadRequest, gin.H{"error": "user_id and day are required"})
		return
	}

	_, err := config.DB.Exec(
		"DELETE FROM planner WHERE user_id=$1 AND day=$2",
		data.UserID, normalizedDay,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to remove plan"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Removed"})
}
