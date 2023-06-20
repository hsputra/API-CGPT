package handlers

import (
	"os"

	"github.com/gin-gonic/gin"
	"github.com/hsputra/API-CGPT/utils"
)

// create Request struct with AdminKey and UserId as fields
type Request struct {
	AdminKey string `json:"admin_key"`
	UserId   string `json:"user_id"`
}

// create function Admin_userAdd with gin.Context as parameter that returns user_id and token
func Admin_userAdd(c *gin.Context) {
	// get admin key from request body
	var request Request
	err := c.BindJSON(&request)
	if err != nil {
		// error 400 Invalid request body
		c.JSON(400, gin.H{
			"error": "Invalid request body",
		})
		return
	}

	// check if utils admin key is valid
	if !utils.VerifyAdminKey(request.AdminKey) {
		// error 401 Unauthorized
		c.JSON(401, gin.H{
			"error": "Invalid admin key",
		})
		return
	}

	// generate user_id and token
	user_id := utils.GenerateId()
	token := utils.GenerateId()

	// insert user_id and token to database
	err = utils.DatabaseInsert(user_id, token)
	if err != nil {
		// error 500 Internal server error
		c.JSON(500, gin.H{
			"error": "Failed to insert user_id and token to database",
		})
		return
	}

	// return user_id and token
	c.JSON(200, gin.H{
		"user_id": user_id,
		"token":   token,
	})
}

// create POST request for admin to delete a user
func Admin_userDelete(c *gin.Context) {
	// get admin key from request body
	var request Request
	if err := c.ShouldBindJSON(&request); err != nil {
		// error 400 Invalid request body
		c.JSON(400, gin.H{
			"error": "Invalid request body",
		})
		return
	}

	// check if utils admin key is valid
	if !utils.VerifyAdminKey(request.AdminKey) {
		// error 401 Unauthorized
		c.JSON(401, gin.H{
			"error": "Invalid admin key",
		})
		return
	}

	// delete user from database
	err := utils.DatabaseDelete(request.UserId)
	if err != nil {
		// error 500 Internal server error
		c.JSON(500, gin.H{
			"error": "Failed to delete user from database",
		})
		return
	}

	// return success message
	c.JSON(200, gin.H{
		"message": "User deleted successfully",
	})
}

// create POST request for admin to get all users
func Admin_usersGet(c *gin.Context) {
	// get admin key from GET parameter
	AdminKey := c.Query("admin_key")

	// check if utils admin key is valid
	if !utils.VerifyAdminKey(AdminKey) {
		// error 401 Unauthorized
		c.JSON(401, gin.H{
			"error":   "Invalid admin key",
			"key":     AdminKey,
			"correct": os.Args[2],
		})
		return
	}

	//  get users from database
	users, err := utils.DatabaseSelectAll()
	if err != nil {
		// error 500 Internal server error
		c.JSON(500, gin.H{
			"error": "Failed to get users from database",
		})
		return
	}

	// return users
	c.JSON(200, gin.H{
		"users": users,
	})
}
