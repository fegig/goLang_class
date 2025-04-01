package controllers

import (
	"net/http"
	"web-service/app/services"
	"web-service/database"
	"web-service/functions"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type UserData struct {
	ID        int    `json:"id"`
	UserId    string `json:"userId"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Email     string `json:"email"`
	Password  string `json:"password"`
}

var users = []UserData{
	{ID: 1, UserId: "dnfg9v8", FirstName: "Test", LastName: "User", Email: "test@test.com", Password: "password"},
	{ID: 2, UserId: "dsfg345", FirstName: "Test", LastName: "User", Email: "test2@test.com", Password: "password2"},
	{ID: 3, UserId: "dg35gdr", FirstName: "Test", LastName: "User", Email: "test3@test.com", Password: "password3"},
}

func GetUsers(c *gin.Context) {
	pagination := functions.GetPagination(c)
	users, err := services.GetUsers(pagination.GetLimit(), pagination.GetOffset())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	functions.SendPaginatedResponse(c, users, pagination)
}

func GetUser(c *gin.Context) {
	userID := c.Param("userID")

	rows, err := database.SelectData("users", []string{"id", "userId", "firstName", "lastName", "email"}, []database.FieldValue{
		{Field: "userId", Value: userID},
	}, nil, nil, nil)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer rows.Close()

	var user UserData
	if rows.Next() {
		if err := rows.Scan(&user.UserId); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	}

	c.IndentedJSON(http.StatusOK, user)
}

func LoginUser(c *gin.Context) {
	email := c.PostForm("email")
	password := c.PostForm("password")

	rows, err := database.SelectData("users", []string{"userId"}, []database.FieldValue{
		{Field: "email", Value: email},
		{Field: "password", Value: password},
	}, nil, nil, nil)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer rows.Close()

	var user UserData
	if rows.Next() {
		if err := rows.Scan(&user.UserId); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	}

	if user.UserId == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Login successful", "userId": user.UserId})
}

func RegisterUser(c *gin.Context) {
	var newUser UserData
	if err := c.ShouldBindJSON(&newUser); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	newUser.UserId = uuid.New().String()
	newUser.Password = functions.HashPassword(newUser.Password)

	_, err := database.InsertData("users", []database.FieldValue{
		{Field: "userId", Value: newUser.UserId},
		{Field: "firstName", Value: newUser.FirstName},
		{Field: "lastName", Value: newUser.LastName},
		{Field: "email", Value: newUser.Email},
		{Field: "password", Value: newUser.Password},
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.IndentedJSON(http.StatusCreated, newUser)
}

func UpdateUser(c *gin.Context) {
	userID := c.Param("userID")
	email := c.PostForm("email")
	password := c.PostForm("password")
	firstName := c.PostForm("firstName")
	lastName := c.PostForm("lastName")

	for i, u := range users {
		if u.UserId == userID {
			users[i] = UserData{ID: u.ID, UserId: u.UserId, FirstName: firstName, LastName: lastName, Email: email, Password: password}
			c.JSON(http.StatusOK, gin.H{"message": "User updated successfully"})
			return
		}
	}
	c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
}
