package handlers

import (
	"database/sql"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/DevSoc-exe/placement-portal-backend/internal/models"
	"github.com/DevSoc-exe/placement-portal-backend/internal/pkg"
	"github.com/DevSoc-exe/placement-portal-backend/internal/responses"

	// "github.com/DevSoc-exe/placement-portal-backend/internal/pkg"
	// "github.com/aidarkhanov/nanoid"
	"github.com/gin-gonic/gin"
	// "golang.org/x/crypto/bcrypt"
)

func FormatTime(t time.Time) string {
	return t.Format("03:04 PM 02/01/2006")
}

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

		// driveBody.DateOfDrive = driveBody.DateOfDrive[0:10]
		var drive models.Drive
		drive.DateOfDrive, drive.Deadline, err = pkg.ParseDates(driveBody.DateOfDrive, driveBody.Deadline)
		if err != nil {
			respError.Message = string("Error parsing dates")
			respError.MapApiResponse(c, http.StatusBadRequest)
			return
		}

		allowed_branches := driveBody.AllowedBranches
		drive.Cse_allowed = strings.Contains(allowed_branches, "Computer Science and Engineering")
		drive.Ece_allowed = strings.Contains(allowed_branches, "Electronics and Communication Engineering")
		drive.Mech_allowed = strings.Contains(allowed_branches, "Mechanical Engineering")
		drive.Civ_allowed = strings.Contains(allowed_branches, "Civil Engineering")

		drive.CompanyID = driveBody.CompanyID
		drive.DriveDuration = driveBody.DriveDuration
		drive.Roles = driveBody.Roles
		drive.Location = driveBody.Location
		drive.Qualifications = driveBody.Qualifications
		drive.PointsToNote = driveBody.PointsToNote
		drive.JobDescription = driveBody.JobDescription
		drive.MinCGPA = driveBody.MinCGPA
		drive.DriveType = driveBody.DriveType
		drive.RequiredData = driveBody.RequiredData

		allowedBranches := strings.Split(allowed_branches, ",")
		mailingList, err := s.GetUserMailsByBranchesAboveCGPA(allowedBranches, drive.MinCGPA)
		company, err := s.GetCompanyUsingCompanyID(driveBody.CompanyID)

		driveID, err := s.CreateNewDriveUsingObject(drive)

		driveCrux := pkg.CompanyCrux{
			Name:     company.Name,
			Deadline: FormatTime(drive.Deadline),
			ID:       driveID,
		}

		// fmt.Println(driveCrux)
		if err != nil {
			respError.Message = string(err.Error())
			respError.MapApiResponse(c, http.StatusInternalServerError)
			return
		}
		mail := pkg.CreateDriveUpdateNotificationEmail(mailingList, driveCrux)
		err = mail.SendEmail()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to Send OTP Email.", "message": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{})

		respSuccess := responses.ApiResponse{
			Success: true,
			Message: string(responses.DriveCreated),
			Data:    nil,
		}
		respSuccess.MapApiResponse(c, http.StatusCreated)
	}
}

func HandleGetDriveUsingID(s models.Store) gin.HandlerFunc {
	return func(c *gin.Context) {
		respError := responses.ApiResponse{
			Success: false,
			Message: "",
			Data:    nil,
		}

		userID, exists := c.Get("userID")
		if !exists {
			respError.Message = string(responses.UserNotFound)
			respError.MapApiResponse(c, http.StatusUnauthorized)
			return
		}

		id, exists := c.Params.Get("id")
		if !exists {
			respError.Message = string(responses.DriveNotFound)
			respError.MapApiResponse(c, http.StatusNotFound)
			return
		}

		data, err := s.GetJobPostingUsingDriveID(id)
		if err != nil {
			if err == sql.ErrNoRows {
				respError.Message = string(responses.DriveNotFound)
				respError.MapApiResponse(c, http.StatusNotFound)
			} else {
				respError.Message = string(err.Error())
				respError.MapApiResponse(c, http.StatusInternalServerError)
			}
			return
		}

		appliedRole, err := s.GetAppliedRole(userID.(string), id)
		if err != nil && err != sql.ErrNoRows {
			respError.Message = string(responses.DatabaseError)
			respError.MapApiResponse(c, http.StatusInternalServerError)
			return
		}

		if appliedRole != nil {
			data.AppliedRole.ID = appliedRole.ID
			data.AppliedRole.DriveID = appliedRole.DriveID
			data.AppliedRole.SalaryHigh = appliedRole.SalaryHigh
			data.AppliedRole.SalaryLow = appliedRole.SalaryLow
			data.AppliedRole.StipendHigh = appliedRole.StipendHigh
			data.AppliedRole.StipendLow = appliedRole.StipendLow
		}

		respSuccess := responses.ApiResponse{
			Success: true,
			Message: string(responses.DriveFound),
			Data:    data,
		}
		respSuccess.MapApiResponse(c, http.StatusOK)
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

func HandleGetCompanyFromID(s models.Store) gin.HandlerFunc {
	return func(c *gin.Context) {

		companyID := c.Query("id")

		company, err := s.GetCompanyUsingCompanyID(companyID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error":   "Internal server error",
				"message": err.Error(),
			})
			return
		}
		respSuccess := responses.ApiResponse{
			Success: true,
			Message: string(responses.CompanyCreated),
			Data:    company,
		}
		respSuccess.MapApiResponse(c, http.StatusCreated)
	}
}

func HandleGetCompaniesForUser(s models.Store) gin.HandlerFunc {
	return func(c *gin.Context) {

		page := c.Query("page")
		if page == "" {
			page = "0"
		}

		name := c.Query("q")

		companies, err := s.GetAllCompaniesForUser(page, name)
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

func HandleGetDrivesForUser(s models.Store) gin.HandlerFunc {
	return func(c *gin.Context) {
		drives, err := s.GetAllDrivesForUser()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error":   "Internal server error",
				"message": err.Error(),
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"drives":       drives,
			"total_drives": len(drives),
		})
	}
}
