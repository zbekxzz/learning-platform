package modules

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(rg *gin.RouterGroup) {

	group := rg.Group("/modules")

	group.POST("/", createModule)
	group.GET("/course/:course_id", getModules)

	group.POST("/material", createMaterial)
	group.GET("/:module_id/materials", getMaterials)
}

func createModule(c *gin.Context) {

	var req struct {
		CourseID   int64  `json:"course_id"`
		Title      string `json:"title"`
		OrderIndex int    `json:"order_index"`
	}

	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid"})
		return
	}

	role := c.GetString("role")

	module, err := CreateModuleService(req.CourseID, req.Title, req.OrderIndex, role)
	if err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, module)
}

func getModules(c *gin.Context) {

	idParam := c.Param("course_id")
	courseID, _ := strconv.ParseInt(idParam, 10, 64)

	modules, err := GetModulesByCourse(courseID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "cannot fetch"})
		return
	}

	c.JSON(http.StatusOK, modules)
}

func createMaterial(c *gin.Context) {

	var m Material

	if err := c.BindJSON(&m); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid"})
		return
	}

	role := c.GetString("role")

	err := CreateMaterialService(&m, role)
	if err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, m)
}

func getMaterials(c *gin.Context) {

	idParam := c.Param("module_id")
	moduleID, _ := strconv.ParseInt(idParam, 10, 64)

	materials, err := GetMaterialsByModule(moduleID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "cannot fetch"})
		return
	}

	c.JSON(http.StatusOK, materials)
}
