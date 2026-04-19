package modules

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"platform/internal/storage"
)

func RegisterRoutes(rg *gin.RouterGroup) {

	group := rg.Group("/modules")

	group.POST("/", createModule)
	group.GET("/chapter/:chapter_id", getModules)

	group.POST("/material", createMaterial)
	group.GET("/:module_id/materials", getMaterials)

	group.POST("/upload", uploadMaterialFile)
}

func createModule(c *gin.Context) {

	var req struct {
		ChapterID  int64  `json:"chapter_id"`
		Title      string `json:"title"`
		OrderIndex int    `json:"order_index"`
	}

	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid"})
		return
	}

	role := c.GetString("role")

	module, err := CreateModuleService(req.ChapterID, req.Title, req.OrderIndex, role)
	if err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, module)
}

func getModules(c *gin.Context) {

	idParam := c.Param("chapter_id")
	chapterID, _ := strconv.ParseInt(idParam, 10, 64)

	modules, err := GetModulesByChapter(chapterID)
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

func uploadMaterialFile(c *gin.Context) {

	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(400, gin.H{"error": "file required"})
		return
	}

	dst := "./tmp/" + file.Filename
	c.SaveUploadedFile(file, dst)

	url, err := storage.UploadFile("materials", file.Filename, dst)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{
		"url": url,
	})
}
