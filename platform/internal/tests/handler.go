package tests

import (
	"github.com/gin-gonic/gin"
	"strconv"
)

func RegisterRoutes(rg *gin.RouterGroup) {

	group := rg.Group("/tests")

	group.GET("/module/:module_id/start", startTest)
	group.POST("/submit", submitTest)
}

func startTest(c *gin.Context) {

	moduleID, _ := strconv.ParseInt(c.Param("module_id"), 10, 64)
	userID := c.GetInt64("user_id")

	test, questions, answers, err := StartTest(userID, moduleID)
	if err != nil {
		c.JSON(500, gin.H{"error": "failed"})
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
		TestID  int64           `json:"test_id"`
		Answers map[int64]int64 `json:"answers"`
	}

	c.BindJSON(&req)

	userID := c.GetInt64("user_id")

	score, _ := SubmitTest(req.Answers, req.TestID, userID)

	c.JSON(200, gin.H{
		"score": score,
	})
}
