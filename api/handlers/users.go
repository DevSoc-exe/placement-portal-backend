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
		name := c.Query("q")
		branch := c.Query("branch")
		role := c.Query("role")


		students, err := s.GetAllStudents(page, gender, branch, role, name)
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

func HandleToggleUserRole(s models.Store) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		if id == "" {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "ID Not Found",
			})
			return
		}

		err := s.ToggleUserRole(id)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error":   "Internal server error",
				"message": err.Error(),
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"message": "User role toggled successfully",
		})
	}
}
