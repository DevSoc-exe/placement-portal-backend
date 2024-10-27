package database

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/DevSoc-exe/placement-portal-backend/internal/models"
	_ "github.com/DevSoc-exe/placement-portal-backend/internal/models"
	_ "github.com/DevSoc-exe/placement-portal-backend/internal/responses"
	_ "github.com/go-sql-driver/mysql"
	"github.com/google/uuid"
)

func (db *Database) createJobsTable() error {

	companyTableQuery := `CREATE TABLE IF NOT EXISTS company (
		company_id VARCHAR(36) PRIMARY KEY NOT NULL,
		name VARCHAR(255) NOT NULL,
		overview LONGTEXT NOT NULL,
		contact_email VARCHAR(255) NOT NULL,
		contact_number VARCHAR(20),
		linked_in VARCHAR(255),
		website VARCHAR(255)
	);`

	driveTableQuery := `CREATE TABLE IF NOT EXISTS drive (
		id VARCHAR(36) PRIMARY KEY NOT NULL,
		company_id VARCHAR(36),
		drive_date DATE NOT NULL,
		drive_duration INT NOT NULL,
		location VARCHAR(255),
		key_responsibilities LONGTEXT NOT NULL,
		qualifications LONGTEXT NOT NULL,
		points_to_note LONGTEXT NOT NULL,
		job_description LONGBLOB,
		FOREIGN KEY (company_id) REFERENCES company(company_id) ON DELETE CASCADE
	);`

	roleTableQuery := `CREATE TABLE IF NOT EXISTS role (
		id VARCHAR(36) PRIMARY KEY NOT NULL,
		drive_id VARCHAR(36),
		title VARCHAR(255) NOT NULL,
		stipend_low INT,
		stipend_high INT,
		salary_low INT,
		salary_high INT,
		FOREIGN KEY (drive_id) REFERENCES drive(id) ON DELETE CASCADE
	);`

	_, err := db.DB.Exec(companyTableQuery)
	if err != nil {
		return fmt.Errorf("could not create company table: %s", err)
	}

	_, err = db.DB.Exec(driveTableQuery)
	if err != nil {
		return fmt.Errorf("could not create drive table: %s", err)
	}

	_, err = db.DB.Exec(roleTableQuery)
	if err != nil {
		return fmt.Errorf("could not create role table: %s", err)
	}

	return nil
}

func (db *Database) CreateNewDriveUsingObject(driveData models.DriveBody) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*2)
	defer cancel()

	queryToInsertDrive := `
	INSERT INTO drive (id, company_id, drive_date, drive_duration, location, key_responsibilities, qualifications, points_to_note, job_description)
	VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?);
	`

	fmt.Println(driveData.CompanyID)
	date, err := time.Parse("2006-01-02", driveData.DateOfDrive)
	_, err = db.DB.ExecContext(ctx, queryToInsertDrive, driveData.ID, driveData.CompanyID, date, driveData.DriveDuration, driveData.Location, driveData.Responsibilities, driveData.Qualifications, driveData.PointsToNote, driveData.JobDescription)
	if err != nil {
		fmt.Println("error was here!")
		return err
	}

	queryToInsertRoles := `
	INSERT INTO role (id, drive_id, title, stipend_low, stipend_high, salary_low, salary_high)
	VALUES `

	var valueStrings []string
	var valueArgs []interface{}

	for _, role := range driveData.Roles {
		roleUUID := uuid.New().String()
		valueStrings = append(valueStrings, "(?, ?, ?, ?, ?, ?, ?)")
		valueArgs = append(valueArgs, roleUUID, role.DriveID, role.Title, role.StipendLow, role.StipendHigh, role.SalaryLow, role.SalaryHigh)
	}

	queryToInsertRoles += strings.Join(valueStrings, ", ")

	_, err = db.DB.ExecContext(ctx, queryToInsertRoles, valueArgs...)
	fmt.Println(queryToInsertRoles)
	if err != nil {
		fmt.Println("error was here 2!")
		return err
	}

	return err
}

func (db *Database) DeleteJobUsingDriveID(driveID string) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*2)
	defer cancel()

	queryToDeleteRoles := `
	DELETE
	FROM drive
	WHERE id = ?;
	`

	_, err := db.DB.ExecContext(ctx, queryToDeleteRoles, driveID)
	if err != nil {
		if err == sql.ErrNoRows {
			return err
		}
		return err
	}
	return err
}

func (db *Database) GetJobPostingUsingDriveID(driveID string) (interface{}, error) {

	// user_id, exists := c.Get("user_id")
	// if !exists {
	// 	c.AbortWithStatus(http.StatusUnauthorized)
	// 	return

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*2)
	defer cancel()

	queryToGetDriveInfo := `
	SELECT company.company_id, name, overview, contact_email, contact_number, linked_in, website, drive_date, drive_duration, location, key_responsibilities, qualifications, points_to_note, job_description
	FROM company 
	JOIN drive ON company.company_id = drive.company_id
	WHERE drive.id = ?;
	`
	row := db.DB.QueryRowContext(ctx, queryToGetDriveInfo, driveID)
	var drive models.Drive
	err := row.Scan(&drive.CompanyID, &drive.Company.Name, &drive.Company.Overview, &drive.Company.ContactEmail, &drive.Company.ContactNumber, &drive.Company.LinkedIn, &drive.Company.Website, &drive.DateOfDrive, &drive.DriveDuration, &drive.Location, &drive.Responsibilities, &drive.Qualifications, &drive.PointsToNote, &drive.JobDescription)

	if err != nil {
		return drive, err
	}

	roles, err := db.GetRolesUsingDriveID(driveID)

	if err != nil {
		return drive, err
	}
	drive.Roles = roles
	return drive, err
}

func (db *Database) GetRolesUsingDriveID(driveID string) ([]models.Role, error) {
	roles := make([]models.Role, 0)

	queryToGetRoles := `
    SELECT id, drive_id, title, stipend_low, stipend_high, salary_low, salary_high
    FROM role
    WHERE drive_id = ?;
    `

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*2)
	defer cancel()

	rows, err := db.DB.QueryContext(ctx, queryToGetRoles, driveID)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		role := new(models.Role)
		if err := rows.Scan(&role.ID, &role.DriveID, &role.Title, &role.StipendLow, &role.StipendHigh, &role.SalaryLow, &role.SalaryHigh); err != nil {
			return nil, err
		}
		roles = append(roles, *role)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return roles, nil
}
