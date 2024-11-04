package handlers

import (
	"database/sql"
	"fmt"
	"net/http"

	"github.com/DevSoc-exe/placement-portal-backend/internal/models"
	"github.com/DevSoc-exe/placement-portal-backend/internal/responses"

	// "github.com/DevSoc-exe/placement-portal-backend/internal/pkg"
	// "github.com/aidarkhanov/nanoid"
	"github.com/gin-gonic/gin"
	// "golang.org/x/crypto/bcrypt"
)

func HandleCreateNewDrive(s models.Store) gin.HandlerFunc {
	return func(c *gin.Context) {
		// user_id, exists := c.Get("user_id")
		// if !exists {cc
		// 	c.AbortWithStatus(http.StatusUnauthorized)
		// 	return
		// }

		respError := responses.ApiResponse{
			Success: false,
			Message: "",
			Data:    nil,
		}

		var driveBody models.DriveBody

		err := c.BindJSON(&driveBody)
		if err != nil {
			respError.Message = string(responses.BindError)
			respError.MapApiResponse(c, http.StatusBadRequest)
			return
		}

		fmt.Println(driveBody)
		err = s.CreateNewDriveUsingObject(driveBody)

		if err != nil {
			respError.Message = string(err.Error())
			respError.MapApiResponse(c, http.StatusInternalServerError)
			return
		}

		respSuccess := responses.ApiResponse{
			Success: true,
			Message: string(responses.DriveCreated),
			Data:    nil,
		}
		respSuccess.MapApiResponse(c, http.StatusCreated)
		return
	}
}

func HandleGetDriveUsingID(s models.Store) gin.HandlerFunc {
	return func(c *gin.Context) {
		fmt.Print("Inside Backend!")
		var body struct {
			DriveID string
		}

		respError := responses.ApiResponse{
			Success: false,
			Message: "",
			Data:    nil,
		}

		if err := c.Bind(&body); err != nil {
			respError.Message = string(responses.BindError)
			respError.MapApiResponse(c, http.StatusBadRequest)
		}

		data, err := s.GetJobPostingUsingDriveID(body.DriveID)

		fmt.Println(err)
		if err == sql.ErrNoRows {
			respError.Message = string(responses.DatabaseError)
			respError.MapApiResponse(c, http.StatusInternalServerError)
		}

		respSuccess := responses.ApiResponse{
			Success: true,
			Message: string(responses.DriveFound),
			Data:    data,
		}
		respSuccess.MapApiResponse(c, http.StatusFound)
	}
}

func HandleDeleteDrive(s models.Store) gin.HandlerFunc {
	return func(c *gin.Context) {

		// user_id, exists := c.Get("user_id")
		// if !exists {
		// 	c.AbortWithStatus(http.StatusUnauthorized)
		// 	return
		// }

		var body struct {
			DriveID string
		}

		respError := responses.ApiResponse{
			Success: false,
			Message: string(responses.DriveNotFound),
			Data:    nil,
		}

		if err := c.Bind(&body); err != nil {
			respError.MapApiResponse(c, http.StatusBadRequest)
			return

		}
		driveToDelete := body.DriveID
		fmt.Println(driveToDelete)

		data, err := s.GetJobPostingUsingDriveID(driveToDelete)

		if err != nil {
			respError.Message = err.Error()
			respError.MapApiResponse(c, http.StatusNotFound)
			return
		}

		err = s.DeleteJobUsingDriveID(driveToDelete)
		if err != nil {
			respError.Message = err.Error()
			respError.MapApiResponse(c, http.StatusNotFound)
			return

		}
		respSuccess := responses.ApiResponse{
			Success: true,
			Message: string(responses.DriveFound),
			Data:    data,
		}
		respSuccess.MapApiResponse(c, http.StatusFound)
	}
}

func HandleCreateNewCompany(s models.Store) gin.HandlerFunc {
	return func(c *gin.Context) {
		// user_id, exists := c.Get("user_id")
		// if !exists {
		// 	c.AbortWithStatus(http.StatusUnauthorized)
		// 	return
		// }

		respError := responses.ApiResponse{
			Success: false,
			Message: "",
			Data:    nil,
		}

		var company models.Company

		err := c.BindJSON(&company)
		if err != nil {
			respError.Message = string(responses.BindError)
			respError.MapApiResponse(c, http.StatusBadRequest)
			return
		}

		err = s.AddNewCompany(&company)

		if err != nil {
			respError.Message = string(err.Error())
			respError.MapApiResponse(c, http.StatusInternalServerError)
			return
		}

		respSuccess := responses.ApiResponse{
			Success: true,
			Message: string(responses.CompanyCreated),
			Data:    nil,
		}
		respSuccess.MapApiResponse(c, http.StatusCreated)
	}
}

func HandleGetAllCompanies(s models.Store) gin.HandlerFunc {
	return func(c *gin.Context) {

		page := c.Query("page")
		if page == "" {
			page = "0"
		}

		name := c.Query("q")

		companies, err := s.GetAllCompanies(page, name)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error":   "Internal server error",
				"message": err.Error(),
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"companies":       companies,
			"total_companies": len(companies),
		})
	}
}
