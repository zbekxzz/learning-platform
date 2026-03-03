package enrollments

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(rg *gin.RouterGroup) {

	group := rg.Group("/enrollments")

	group.POST("/:course_id", enroll)
	group.DELETE("/:course_id", unenroll)
	group.GET("/my", myEnrollments)
}

func enroll(c *gin.Context) {

	userID := c.GetInt64("user_id")
	role := c.GetString("role")

	courseParam := c.Param("course_id")
	courseID, err := strconv.ParseInt(courseParam, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid course id"})
		return
	}

	err = EnrollUser(userID, courseID, role)
	if err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "enrolled successfully"})
}

func unenroll(c *gin.Context) {

	userID := c.GetInt64("user_id")
	role := c.GetString("role")

	courseParam := c.Param("course_id")
	courseID, err := strconv.ParseInt(courseParam, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid course id"})
		return
	}

	err = UnenrollUser(userID, courseID, role)
	if err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "unenrolled successfully"})
}

func myEnrollments(c *gin.Context) {

	userID := c.GetInt64("user_id")

	enrollments, err := GetMyEnrollments(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "cannot fetch enrollments"})
		return
	}

	c.JSON(http.StatusOK, enrollments)
}
