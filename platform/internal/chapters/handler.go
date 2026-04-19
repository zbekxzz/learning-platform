package chapters

import (
	"github.com/gin-gonic/gin"
	"strconv"
)

func RegisterRoutes(rg *gin.RouterGroup) {

	group := rg.Group("/chapters")

	group.POST("/", createChapter)
	group.GET("/course/:course_id", getChapters)
}

func createChapter(c *gin.Context) {

	var req struct {
		CourseID   int64  `json:"course_id"`
		Title      string `json:"title"`
		OrderIndex int    `json:"order_index"`
	}

	c.BindJSON(&req)

	role := c.GetString("role")

	if role != "admin" && role != "teacher" {
		c.JSON(403, gin.H{"error": "forbidden"})
		return
	}

	ch := &Chapter{
		CourseID:   req.CourseID,
		Title:      req.Title,
		OrderIndex: req.OrderIndex,
	}

	CreateChapter(ch)

	c.JSON(200, ch)
}

func getChapters(c *gin.Context) {

	courseID, _ := strconv.ParseInt(c.Param("course_id"), 10, 64)

	ch, _ := GetChaptersByCourse(courseID)

	c.JSON(200, ch)
}
