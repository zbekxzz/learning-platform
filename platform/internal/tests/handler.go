package tests

import (
	"github.com/gin-gonic/gin"
	"strconv"
)

func RegisterRoutes(rg *gin.RouterGroup) {

	group := rg.Group("/tests")

	group.GET("/module/:module_id/start", startTest)
	group.GET("/chapter/:chapter_id/start", startChapterTest)
	group.GET("/course/:course_id/final/start", startFinalTest)
	group.POST("/submit", submitTest)
	group.POST("/create-full", createFullTest)
}

func startTest(c *gin.Context) {

	moduleID, _ := strconv.ParseInt(c.Param("module_id"), 10, 64)
	userID := c.GetInt64("user_id")

	test, questions, answers, err := StartTest(userID, moduleID)
	if err != nil {

		if err.Error() == "module locked" {
			c.JSON(403, gin.H{
				"error": "Алдымен алдыңғы модульді аяқтаңыз",
			})
			return
		}

		c.JSON(500, gin.H{"error": "failed"})
		return
	}

	c.JSON(200, gin.H{
		"test":      test,
		"questions": questions,
		"answers":   answers,
	})
}

func startChapterTest(c *gin.Context) {

	chapterID, _ := strconv.ParseInt(c.Param("chapter_id"), 10, 64)
	userID := c.GetInt64("user_id")

	test, questions, answers, err := StartChapterTest(userID, chapterID)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{
		"test":      test,
		"questions": questions,
		"answers":   answers,
	})
}

func startFinalTest(c *gin.Context) {

	courseID, _ := strconv.ParseInt(c.Param("course_id"), 10, 64)
	userID := c.GetInt64("user_id")

	test, questions, answers, err := StartFinalTest(userID, courseID)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{
		"test":      test,
		"questions": questions,
		"answers":   answers,
	})
}

func submitTest(c *gin.Context) {
	var req struct {
		TestID  int64                  `json:"test_id"`
		Answers map[string]interface{} `json:"answers"`
	}

	if err := c.BindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": "invalid request"})
		return
	}

	userID := c.GetInt64("user_id")

	score, err := SubmitTest(req.Answers, req.TestID, userID)
	if err != nil {

		if err.Error() == "test already passed" {
			c.JSON(400, gin.H{
				"error": "Тест уже пройден",
			})
			return
		}

		if err.Error() == "no attempts left" {
			c.JSON(400, gin.H{
				"error": "Попытки закончились",
			})
			return
		}

		c.JSON(500, gin.H{
			"error": err.Error(),
		})
		return
	}
	c.JSON(200, gin.H{
		"score": score,
	})
}

func createFullTest(c *gin.Context) {

	var req CreateTestRequest

	if err := c.BindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": "invalid request"})
		return
	}

	role := c.GetString("role")

	if role != "admin" && role != "teacher" {
		c.JSON(403, gin.H{"error": "forbidden"})
		return
	}

	err := CreateFullTest(req)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{"message": "test created"})
}
