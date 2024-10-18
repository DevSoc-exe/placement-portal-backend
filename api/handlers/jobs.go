package handlers

import (
	"context"
	// "fmt"
	"net/http"
	"strings"

	// "fmt"
	// "net/http"
	// "database/sql"
	// "os"
	"time"

	"github.com/DevSoc-exe/placement-portal-backend/internal/models"
	"github.com/DevSoc-exe/placement-portal-backend/internal/responses"

	// "github.com/DevSoc-exe/placement-portal-backend/internal/pkg"
	// "github.com/aidarkhanov/nanoid"
	"github.com/gin-gonic/gin"
	// "golang.org/x/crypto/bcrypt"
)

func (db *Database) CreateJobPosting(c *gin.Context) error {
	// user_id, exists := c.Get("user_id")
	// if !exists {
	// 	c.AbortWithStatus(http.StatusUnauthorized)
	// 	return
	// }

	var body struct {
		DriveID          string
		CompanyID        string
		DateOfDrive      time.Time
		DriveDuration    int
		Roles            []models.Role
		Location         string
		Responsibilities string
		Qualifications   string
		PointsToNote     string
		JobDescription   []byte
	}

	err := c.BindJSON(&body)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": responses.IvalidJobPosting,
		})
		return err
	}

	// current_time := time.Now().UTC()

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*2)
	defer cancel()

	queryToInsertRoles := `
	INSERT INTO roles (id, drive_id, title, stipend_low, stipend_high, salary_low, salary_high, created_at, updated_at)
	VALUES
	`

	var valueStrings []string
	var valueArgs []interface{}

	for _, role := range body.Roles {
		valueStrings = append(valueStrings, "(?, ?, ?, ?, ?, ?, NOW(), NOW())")
		valueArgs = append(valueArgs, role.ID, role.DriveID, role.Title, role.StipendLow, role.StipendHigh, role.SalaryLow, role.SalaryHigh)
	}

	queryToInsertRoles += strings.Join(valueStrings, ", ")

	_, err = db.DB.ExecContext(ctx, queryToInsertRoles, valueArgs...)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": string(responses.RolesInsertionFail) + err.Error(),
		})
		return err
	}

	queryToInsertDrive := `
    INSERT INTO drive (id, company_id, drive_date, drive_duration, location, key_responsibilities, qualifications, points_to_note,job_description, created_at, updated_at)
    VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, NOW(), NOW());
    `

	_, err = db.DB.ExecContext(ctx, queryToInsertDrive, body.DriveID, body.CompanyID, body.DateOfDrive, body.DriveDuration, body.Location, body.Responsibilities, body.Qualifications, body.PointsToNote, body.JobDescription)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": string(responses.JobPostingFailed) + err.Error(),
		})
		return err
	}

	if err != nil {
		return err
	}

	return nil
}
