package progress

import (
	"net/http"
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(rg *gin.RouterGroup) {
	group := rg.Group("/progress")
	group.POST("/complete", completeModule)
	group.GET("/completed", getCompleted)
}

func completeModule(c *gin.Context) {
	var req struct {
		ModuleID int64 `json:"module_id"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid payload"})
		return
	}

	userID := c.GetInt64("user_id")

	err := MarkCompleted(userID, req.ModuleID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "completed"})
}

func getCompleted(c *gin.Context) {
	userID := c.GetInt64("user_id")

	completedList, err := GetCompletedModules(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if completedList == nil {
		completedList = []int64{}
	}

	c.JSON(http.StatusOK, completedList)
}
