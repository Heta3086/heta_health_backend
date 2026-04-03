package handlers

import (
	"net/http"
	"strings"

	"heta_health_backend/config"

	"github.com/gin-gonic/gin"
)

func GetDietOptions(c *gin.Context) {
	rows, err := config.DB.Query(`
		SELECT DISTINCT CAST(diet_type AS TEXT)
		FROM meals
		WHERE diet_type IS NOT NULL
		ORDER BY CAST(diet_type AS TEXT)
	`)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch diet options", "detail": err.Error()})
		return
	}
	defer rows.Close()

	options := make([]gin.H, 0)
	for rows.Next() {
		var value string
		if err := rows.Scan(&value); err != nil {
			continue
		}

		clean := strings.TrimSpace(value)
		if clean == "" {
			continue
		}

		options = append(options, gin.H{
			"value": clean,
			"label": formatDietLabel(clean),
		})
	}

	c.JSON(http.StatusOK, options)
}

func formatDietLabel(value string) string {
	switch strings.ToLower(value) {
	case "veg":
		return "Vegetarian"
	case "nonveg", "non_veg":
		return "Non Veg"
	case "eggetarian":
		return "Eggetarian"
	case "vegan":
		return "Vegan"
	default:
		formatted := strings.ReplaceAll(strings.ToLower(value), "_", " ")
		if formatted == "" {
			return "Unknown"
		}
		return strings.ToUpper(formatted[:1]) + formatted[1:]
	}
}
