package main

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

type user struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

var users = []user{
	{ID: "1", Name: "John Doe"},
}

func home(c *gin.Context) {
	example := c.MustGet("example").(string)
	log.Println(example)
	c.String(http.StatusOK, "Welcome to REST API")
}

func getAllUsers(c *gin.Context) {
	c.JSON(http.StatusOK, users)
}

func main() {
	router := gin.Default()
	router.Use(Logger())
	router.GET("/", home)
	router.GET("/users", getAllUsers)
	router.Run("localhost:8000")
}
