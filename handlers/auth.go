package handlers

import (
	"heta_health_backend/config"
	"net/http"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

func Signup(c *gin.Context) {
	var user struct {
		Name     string `json:"name"`
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	c.BindJSON(&user)

	if user.Email == "" || user.Password == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Email and password are required"})
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to secure password"})
		return
	}

	_, err = config.DB.Exec(
		"INSERT INTO auth_users (name, email, password) VALUES ($1,$2,$3)",
		user.Name, user.Email, string(hashedPassword),
	)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "User already exists"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Signup successful"})
}

func Login(c *gin.Context) {
	var user struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	c.BindJSON(&user)

	if user.Email == "" || user.Password == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Email and password are required"})
		return
	}

	var id int
	var name string
	var dbPassword string

	err := config.DB.QueryRow(
		"SELECT id, name, password FROM auth_users WHERE email=$1",
		user.Email,
	).Scan(&id, &name, &dbPassword)

	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(dbPassword), []byte(user.Password))
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Login successful",
		"user_id": id,
		"name":    name,
	})
}

func Logout(c *gin.Context) {
	var input struct {
		UserID int `json:"user_id"`
	}

	_ = c.BindJSON(&input)

	c.JSON(http.StatusOK, gin.H{
		"message": "Logout successful",
	})
}
