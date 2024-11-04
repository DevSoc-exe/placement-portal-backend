package handlers

import (
	"net/http"

	"github.com/DevSoc-exe/placement-portal-backend/internal/models"
	"github.com/gin-gonic/gin"
)

func HandleGetAllStudents(s models.Store) gin.HandlerFunc {
	return func(c *gin.Context) {
		page := c.Query("page")
		if page == "" {
			page = "0"
		}

		gender := c.Query("gender")
		branch := c.Query("branch")

		students, err := s.GetAllStudents(page, gender, branch)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error":   "Internal server error",
				"message": err.Error(),
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"users":       students,
			"total_users": len(students),
		})
	}
}
