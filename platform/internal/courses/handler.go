package courses

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func RegisterPublicRoutes(rg *gin.RouterGroup) {
	group := rg.Group("/courses")
	group.GET("/", getPublicCourses)
	group.GET("/:id", getPublicCourseByID)
}

func RegisterProtectedRoutes(rg *gin.RouterGroup) {
	group := rg.Group("/courses")

	group.POST("/", create)
	group.DELETE("/:id", deleteCourse)
	group.PUT("/:id/publish", publish)
	group.GET("/admin/all", getAllForAdmin)
}

func getPublicCourses(c *gin.Context) {

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	offset := (page - 1) * limit

	courses, total, err := GetPublishedPaginated(limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "cannot fetch"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data":  courses,
		"total": total,
		"page":  page,
		"limit": limit,
	})
}

func getAllForAdmin(c *gin.Context) {

	role := c.GetString("role")

	if role != "admin" {
		c.JSON(http.StatusForbidden, gin.H{"error": "forbidden"})
		return
	}

	courses, err := GetAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "cannot fetch"})
		return
	}

	c.JSON(http.StatusOK, courses)
}

func getPublicCourseByID(c *gin.Context) {

	idParam := c.Param("id")
	id, _ := strconv.ParseInt(idParam, 10, 64)

	course, err := GetByID(id)
	if err != nil || !course.IsPublished {
		c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
		return
	}

	c.JSON(http.StatusOK, course)
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

func deleteCourse(c *gin.Context) {

	idParam := c.Param("id")
	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	role := c.GetString("role")

	err = DeleteCourse(id, role)
	if err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "course deleted"})
}
