package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

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
	// LoggerWithFormatter middleware will write the logs to gin.DefaultWriter
	// By default gin.DefaultWriter = os.Stdout
	router.Use(gin.LoggerWithFormatter(func(param gin.LogFormatterParams) string {
		// your custom format
		return fmt.Sprintf("%s - [%s] \"%s %s %s %d %s \"%s\" %s\"\n",
			param.ClientIP,
			param.TimeStamp.Format(time.RFC1123),
			param.Method,
			param.Path,
			param.Request.Proto,
			param.StatusCode,
			param.Latency,
			param.Request.UserAgent(),
			param.ErrorMessage,
		)
	}))
	router.Use(gin.Recovery())
	router.Use(Logger())
	router.GET("/", home)
	router.GET("/users", getAllUsers)
	router.Run("localhost:8000")
}
