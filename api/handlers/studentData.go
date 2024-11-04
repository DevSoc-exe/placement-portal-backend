package handlers

import (
	"net/http"

	"github.com/DevSoc-exe/placement-portal-backend/internal/models"
	"github.com/gin-gonic/gin"
)

func HandleAddNewStudentData(s models.Store) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, exists := c.Get("userID")
		if !exists {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "ID Not Found",
			})
			return

		}

		var studentData models.StudentData
		studentData.ID = id.(string)

		err := c.BindJSON(&studentData)
		if err != nil {
			c.JSON(400, gin.H{
				"error": "Invalid request",
			})
			return
		}

		err = s.AddStudentData(&studentData)
		if err != nil {
			c.JSON(500, gin.H{
				"error": "Internal server error",
			})
			return
		}

		c.JSON(201, gin.H{
			"message": "Student data added successfully",
		})
	}
}

func HandleGetStudentData(s models.Store) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, exists := c.Get("userID")
		if !exists {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "ID Not Found",
			})
			return
		}

		studentData, err := s.GetStudentDataByID(id.(string))
		if err != nil {
			c.JSON(500, gin.H{
				"error": "Internal server error",
			})
			return
		}

		c.JSON(200, studentData)
	}
}

func HandleUpdateStudentData(s models.Store) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, exists := c.Get("userID")
		if !exists {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "ID Not Found",
			})
			return
		}

		var studentData models.StudentData
		studentData.ID = id.(string)

		err := c.BindJSON(&studentData)
		if err != nil {
			c.JSON(400, gin.H{
				"error": "Invalid request",
			})
			return
		}

		err = s.UpdateStudentData(&studentData)
		if err != nil {
			c.JSON(500, gin.H{
				"error": "Internal server error",
			})
			return
		}

		c.JSON(200, gin.H{
			"message": "Student data updated successfully",
		})
	}
}

func HandleDeleteStudentData(s models.Store) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, exists := c.Get("userID")
		if !exists {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "ID Not Found",
			})
			return
		}

		err := s.DeleteStudentData(id.(string))
		if err != nil {
			c.JSON(500, gin.H{
				"error": "Internal server error",
			})
			return
		}

		c.JSON(200, gin.H{
			"message": "Student data deleted successfully",
		})
	}
}

//* Admin Routes

func HandleGetAllStudentData(s models.Store) gin.HandlerFunc {
	return func(c *gin.Context) {
		page := c.Query("page")
		if page == "" {
			page = "1"
		}

		gender := c.Query("gender")
		if gender == "" {
			gender = "-"
		}

		branch := c.Query("branch")
		if branch == "" {
			branch = "-"
		}

		studentData, err := s.GetAllStudentData(page, gender, branch)
		if err != nil {
			c.JSON(500, gin.H{
				"error": "Internal server error",
			})
			return
		}

		c.JSON(200, studentData)
	}
}

func HandleGetStudentDataByID(s models.Store) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")

		studentData, err := s.GetStudentDataByID(id)
		if err != nil {
			c.JSON(500, gin.H{
				"error": "Internal server error",
			})
			return
		}

		c.JSON(200, studentData)
	}
}
