package courses

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"platform/internal/chapters"
	"platform/internal/modules"
	"platform/internal/tests"
)

func GetCourseStructure(c *gin.Context) {

	courseID, _ := strconv.ParseInt(c.Param("id"), 10, 64)

	// 1. главы
	chList, err := chapters.GetChaptersByCourse(courseID)
	if err != nil {
		c.JSON(500, gin.H{"error": "failed to get chapters"})
		return
	}

	result := []gin.H{}

	for _, ch := range chList {

		// 2. темы (modules)
		modList, _ := modules.GetModulesByChapter(ch.ID)

		// 3. тест главы
		var chTest interface{} = nil
		t, err := tests.GetTestByChapter(ch.ID)
		if err == nil {
			chTest = t
		}

		result = append(result, gin.H{
			"chapter":      ch,
			"modules":      modList,
			"chapter_test": chTest,
		})
	}

	c.JSON(http.StatusOK, result)
}

