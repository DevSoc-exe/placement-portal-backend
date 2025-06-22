package handlers

import (
	"net/http"

	"github.com/DevSoc-exe/placement-portal-backend/internal/models"
	"github.com/DevSoc-exe/placement-portal-backend/internal/responses"
	"github.com/gin-gonic/gin"
)

func HandleApplyToDrive(s models.Store) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, exists := c.Get("userID")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			return
		}

		var req models.Application
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		err := s.ApplyForDrive(id.(string), req.RoleID, req.DriveID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(201, gin.H{})
	}
}


func HandleMarkStudentAsPlaced(s models.Store) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, exists := c.Get("userID")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			return
		}

		respError := responses.ApiResponse{
			Success: false,
			Message: "",
			Data:    nil,
		}

		applicationId := c.Param("application_id")

		if applicationId == "" {
			respError.Message = "application_id is required"
			c.JSON(http.StatusBadRequest, respError)
			return
		}

		err := s.MarkStudentAsPlaced(id.(string), applicationId)
		if err != nil {
			respError.Message = "could not mark student as placed"
			c.JSON(http.StatusInternalServerError, respError)
			return
		}

		c.Status(http.StatusNoContent)
	}
}
