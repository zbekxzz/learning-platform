package users

import (
	"github.com/gin-gonic/gin"
	"strconv"
)

func RegisterRoutes(rg *gin.RouterGroup) {

	group := rg.Group("/users")

	group.GET("/", getUsers)
	group.POST("/", createUser)
	group.PUT("/:id", updateUser)
	group.DELETE("/:id", deleteUser)
}

func getUsers(c *gin.Context) {

	role := c.GetString("role")

	users, err := GetUsersService(role)
	if err != nil {
		c.JSON(403, gin.H{"error": "forbidden"})
		return
	}

	c.JSON(200, users)
}

func createUser(c *gin.Context) {

	var req struct {
		Email    string `json:"email"`
		FullName string `json:"full_name"`
		Password string `json:"password"`
		Role     string `json:"role"`
	}

	c.BindJSON(&req)

	role := c.GetString("role")

	user, err := CreateUserService(
		req.Email,
		req.FullName,
		req.Password,
		req.Role,
		role,
	)

	if err != nil {
		c.JSON(403, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, user)
}

func updateUser(c *gin.Context) {

	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)

	var req struct {
		Email    string `json:"email"`
		FullName string `json:"full_name"`
		Role     string `json:"role"`
	}

	c.BindJSON(&req)

	role := c.GetString("role")

	err := UpdateUserService(id, req.Email, req.FullName, req.Role, role)
	if err != nil {
		c.JSON(403, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{"message": "updated"})
}

func deleteUser(c *gin.Context) {

	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)

	role := c.GetString("role")

	err := DeleteUserService(id, role)
	if err != nil {
		c.JSON(403, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{"message": "deleted"})
}
