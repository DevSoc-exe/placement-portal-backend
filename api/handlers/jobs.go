package handlers

import (
	"database/sql"
	"encoding/csv"
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
		mailingList, _ := s.GetUserMailsByBranchesAboveCGPA(allowedBranches, drive.MinCGPA)
		company, _ := s.GetCompanyUsingCompanyID(driveBody.CompanyID)

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
				respError.Message = string(responses.DatabaseError)
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

func HandleGetApplicantsForDrive(s models.Store) gin.HandlerFunc {
	return func(c *gin.Context) {
		// TODO: require handling of query params
		driveID := c.Param("id")

		if driveID == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "drive id not found"})
			return
		}

		students, err := s.GetApplicantsForDrive(driveID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error":   "Internal server error",
				"message": err.Error(),
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"success":        true,
			"students":       students,
			"total_students": len(students),
		})
	}
}

type DriveApplicantRequestBody struct {
	RequiredData string `json:"required_data" binding:"required"`
}

func HandleGetDriveApplicantsForRole(s models.Store) gin.HandlerFunc {
	return func(c *gin.Context) {
		roleID := c.Query("rid")
		driveID := c.Query("did")

		if roleID == "" || driveID == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "ids not found in context"})
			return
		}

		var body DriveApplicantRequestBody
		err := c.BindJSON(&body)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "error parsing body"})
			return
		}

		rows, columns, err := s.GetDriveApplicantsForRole(roleID, body.RequiredData, driveID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error":   "Internal server error",
				"message": err.Error(),
			})
			return
		}
		defer rows.Close()

		c.Header("Content-Disposition", "attachment; filename=applicants.csv")
		c.Header("Content-Type", "text/csv")

		writer := csv.NewWriter(c.Writer)
		defer writer.Flush()

		if err := writer.Write(columns); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error":   "Internal server error",
				"message": fmt.Sprintf("failed to write headers to CSV: %v", err),
			})
			return
		}

		values := make([]interface{}, len(columns))
		valuePtrs := make([]interface{}, len(columns))
		for i := range values {
			valuePtrs[i] = &values[i]
		}

		for rows.Next() {
			if err := rows.Scan(valuePtrs...); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{
					"error":   "Internal server error",
					"message": fmt.Sprintf("failed to scan row: %v", err),
				})
				return
			}

			row := make([]string, len(columns))
			for i, val := range values {
				if val == nil {
					row[i] = ""
					continue
				}

				switch v := val.(type) {
				case []byte:
					row[i] = string(v)
				case string:
					row[i] = v
				case int64:
					row[i] = fmt.Sprintf("%d", v)
				case float64:
					row[i] = fmt.Sprintf("%.2f", v)
				case bool:
					row[i] = fmt.Sprintf("%t", v)
				default:
					row[i] = fmt.Sprintf("%v", v)
				}
			}

			if err := writer.Write(row); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{
					"error":   "Internal server error",
					"message": fmt.Sprintf("failed to write row: %v", err),
				})
				return
			}
		}

		if err := rows.Err(); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error":   "Internal server error",
				"message": fmt.Sprintf("rows iteration error: %v", err),
			})
			return
		}
	}
}
