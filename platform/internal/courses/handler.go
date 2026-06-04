package courses

import (
	"context"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"platform/internal/database"
)

func RegisterPublicRoutes(rg *gin.RouterGroup) {
	group := rg.Group("/courses")
	group.GET("/", getPublicCourses)
	group.GET("/:id", getPublicCourseByID)
	group.GET("/:id/structure", GetCourseStructure)
}

func RegisterProtectedRoutes(rg *gin.RouterGroup) {
	group := rg.Group("/courses")

	group.POST("/", create)
	group.DELETE("/:id", deleteCourse)
	group.PUT("/:id/publish", togglePublish)
	group.GET("/admin/all", getAllForAdmin)
	group.GET("/teacher/my", getTeacherCourses)
	group.GET("/:id/statistics", getCourseStatistics)
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

func getTeacherCourses(c *gin.Context) {

	role := c.GetString("role")

	if role != "teacher" && role != "admin" {
		c.JSON(http.StatusForbidden, gin.H{"error": "forbidden"})
		return
	}

	userID := c.GetInt64("user_id")

	courses, err := GetTeacherCourses(userID)
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

func togglePublish(c *gin.Context) {

	idParam := c.Param("id")
	id, _ := strconv.ParseInt(idParam, 10, 64)

	userID := c.GetInt64("user_id")
	role := c.GetString("role")

	err := TogglePublishCourse(id, userID, role)
	if err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "toggled"})
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

func getCourseStatistics(c *gin.Context) {
	courseID, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	userID := c.GetInt64("user_id")
	role := c.GetString("role")

	if role != "teacher" && role != "admin" {
		c.JSON(http.StatusForbidden, gin.H{"error": "forbidden"})
		return
	}

	if role == "teacher" {
		var createdBy int64
		err := database.DB.QueryRow(context.Background(),
			"SELECT created_by FROM courses WHERE id = $1", courseID).Scan(&createdBy)
		if err != nil || createdBy != userID {
			c.JSON(http.StatusForbidden, gin.H{"error": "forbidden"})
			return
		}
	}

	var totalModules int
	err := database.DB.QueryRow(context.Background(),
		`SELECT COUNT(m.id) FROM modules m
		 JOIN chapters c ON c.id = m.chapter_id
		 WHERE c.course_id = $1`, courseID).Scan(&totalModules)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to count modules"})
		return
	}

	rows, err := database.DB.Query(context.Background(),
		`SELECT e.user_id, u.full_name, u.email
		 FROM enrollments e
		 JOIN users u ON u.id = e.user_id
		 WHERE e.course_id = $1`, courseID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch enrollments"})
		return
	}
	defer rows.Close()

	type StudentProgress struct {
		UserID         int64  `json:"user_id"`
		FullName       string `json:"full_name"`
		Email          string `json:"email"`
		CompletedCount int    `json:"completed_count"`
		TotalCount     int    `json:"total_count"`
		Percentage     int    `json:"percentage"`
		CertificateURL string `json:"certificate_url"`
	}

	var stats []StudentProgress

	for rows.Next() {
		var s StudentProgress
		s.TotalCount = totalModules
		err := rows.Scan(&s.UserID, &s.FullName, &s.Email)
		if err != nil {
			continue
		}

		err = database.DB.QueryRow(context.Background(),
			`SELECT COUNT(p.module_id)
			 FROM progress p
			 JOIN modules m ON m.id = p.module_id
			 JOIN chapters c ON c.id = m.chapter_id
			 WHERE p.user_id = $1 AND c.course_id = $2 AND p.is_completed = true`,
			s.UserID, courseID).Scan(&s.CompletedCount)
		if err != nil {
			s.CompletedCount = 0
		}

		if s.TotalCount > 0 {
			s.Percentage = int((float64(s.CompletedCount) / float64(s.TotalCount)) * 100)
		} else {
			s.Percentage = 0
		}

		err = database.DB.QueryRow(context.Background(),
			`SELECT certificate_url FROM certificates WHERE user_id = $1 AND course_id = $2`,
			s.UserID, courseID).Scan(&s.CertificateURL)
		if err != nil {
			s.CertificateURL = ""
		}

		stats = append(stats, s)
	}

	if stats == nil {
		stats = []StudentProgress{}
	}

	c.JSON(http.StatusOK, stats)
}
