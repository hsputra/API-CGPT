package main

import (
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

func main() {

	// get arguments server port and admin key
	// checking if the number of arguments passed to it is less than 3
	if len(os.Args) < 3 {
		//  If so, it prints a usage message with the correct format and exits the program. This is done using the os.Args
		println("Usage: %s <port> <admin-key>\n")
	}

	// print the two arguments passed to it
	println("Port: %s\n", os.Args[1])
	println("Admin Key: %s\n", os.Args[2])

	//Create a database instance and returns an error if the creation process fails by using DatabaseCreate
	// err := utils.DatabaseCreate()
	// if err != nil {
	// 	println("Error creating database: %s\n", err.Error())
	// 	return
	// }

	router := gin.Default()

	// allow CORS for the router
	router.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		// c.Next()
	})

	// add routes to the router
	// register new client connection
	// router.POST("/client/register", handlers.Client_register)

	router.GET("/health", func(c *gin.Context) {
		c.String(http.StatusOK, "OK")
	})

	router.Run(":" + os.Args[1])
}
