package handlers

import (
	"heta_health_backend/config"

	"github.com/gin-gonic/gin"
)

func GetAllergies(c *gin.Context) {
	rows, err := config.DB.Query("SELECT id, name FROM allergies")
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to fetch allergies"})
		return
	}
	defer rows.Close()

	var allergies []map[string]interface{}

	for rows.Next() {
		var id int
		var name string

		err := rows.Scan(&id, &name)
		if err != nil {
			continue
		}

		allergies = append(allergies, gin.H{
			"id": id,
			"name": name,
		})
	}

	c.JSON(200, allergies)
}