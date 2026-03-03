package courses

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(rg *gin.RouterGroup) {

	group := rg.Group("/courses")

	group.GET("/", getAll)
	group.GET("/:id", getByID)

	protected := group.Group("/")
	protected.Use()
	protected.POST("/", create)
	protected.PUT("/:id/publish", publish)
}

func create(c *gin.Context) {

	var req struct {
		Title       string `json:"title"`
		Description string `json:"description"`
	}

	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	userID := c.GetInt64("user_id")
	role := c.GetString("role")

	course, err := CreateCourse(req.Title, req.Description, userID, role)
	if err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, course)
}

func getAll(c *gin.Context) {

	courses, err := GetAllPublished()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "cannot fetch courses"})
		return
	}

	c.JSON(http.StatusOK, courses)
}

func getByID(c *gin.Context) {

	idParam := c.Param("id")
	id, _ := strconv.ParseInt(idParam, 10, 64)

	course, err := GetByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
		return
	}

	c.JSON(http.StatusOK, course)
}

func publish(c *gin.Context) {

	idParam := c.Param("id")
	id, _ := strconv.ParseInt(idParam, 10, 64)

	role := c.GetString("role")

	err := PublishCourse(id, role)
	if err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "published"})
}
