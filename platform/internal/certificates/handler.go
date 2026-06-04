package certificates

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(rg *gin.RouterGroup) {
	group := rg.Group("/certificates")
	group.GET("/my", getMyCertificatesHandler)
}

func getMyCertificatesHandler(c *gin.Context) {
	userID := c.GetInt64("user_id")

	certs, err := GetMyCertificates(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, certs)
}
