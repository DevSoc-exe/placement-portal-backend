package database

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/DevSoc-exe/placement-portal-backend/internal/models"
	"github.com/aidarkhanov/nanoid"
)

func (db *Database) createJobsTable() error {

	companyTableQuery := `CREATE TABLE IF NOT EXISTS company (
		company_id VARCHAR(36) PRIMARY KEY NOT NULL,
		name VARCHAR(255) NOT NULL,
		hr_name VARCHAR(255) NOT NULL,
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
		min_cgpa DECIMAL(3,2) NOT NULL,
		deadline DATETIME NOT NULL,
		qualifications LONGTEXT NOT NULL,
		points_to_note LONGTEXT NOT NULL,
		job_description VARCHAR(255),
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

	tx, err := db.DB.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	queryToInsertDrive := `
	INSERT INTO drive (id, company_id, drive_date, drive_duration, location, qualifications, points_to_note, job_description, min_cgpa, deadline)
	VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?);
	`

	fmt.Println(driveData.CompanyID)
	driveUUID := nanoid.New()
	date, _ := time.Parse("2006-01-02", driveData.DateOfDrive)
	_, err = tx.ExecContext(ctx, queryToInsertDrive, driveUUID, driveData.CompanyID, date, driveData.DriveDuration, driveData.Location, driveData.Qualifications, driveData.PointsToNote, driveData.JobDescription, driveData.MinCGPA, driveData.Deadline)
	if err != nil {
		fmt.Println("error was here!")
		tx.Rollback()
		return err
	}

	queryToInsertRoles := `
	INSERT INTO role (id, drive_id, title, stipend_low, stipend_high, salary_low, salary_high)
	VALUES `

	var valueStrings []string
	var valueArgs []interface{}

	for _, role := range driveData.Roles {
		roleUUID := nanoid.New()
		valueStrings = append(valueStrings, "(?, ?, ?, ?, ?, ?, ?)")
		valueArgs = append(valueArgs, roleUUID, driveUUID, role.Title, role.StipendLow, role.StipendHigh, role.SalaryLow, role.SalaryHigh)
	}

	queryToInsertRoles += strings.Join(valueStrings, ", ")

	_, err = tx.ExecContext(ctx, queryToInsertRoles, valueArgs...)
	fmt.Println(queryToInsertRoles)
	if err != nil {
		tx.Rollback()
		return err
	}

	err = tx.Commit()
	if err != nil {
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

func (db *Database) AddNewCompany(company *models.Company) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*2)
	defer cancel()

	query := `INSERT INTO company (company_id, name, hr_name, overview, contact_email, contact_number, linked_in, website) VALUES (?, ?, ?, ?, ?, ?, ?, ?);`

	// Generate a new Nano ID for the company
	id := nanoid.New()

	_, err := db.DB.ExecContext(ctx, query, id, company.Name, company.HRName, company.Overview, company.ContactEmail, company.ContactNumber, company.LinkedIn, company.Website)
	if err != nil {
		return err
	}
	return nil
}

func (db *Database) GetJobPostingUsingDriveID(driveID string) (interface{}, error) {

	// user_id, exists := c.Get("user_id")
	// if !exists {
	// 	c.AbortWithStatus(http.StatusUnauthorized)
	// 	return

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*2)
	defer cancel()

	queryToGetDriveInfo := `
	SELECT company.company_id, name, overview, contact_email, contact_number, linked_in, website, drive_date, drive_duration, location, qualifications, points_to_note, job_description
	FROM company
	JOIN drive ON company.company_id = drive.company_id
	WHERE drive.id = ?;
	`
	row := db.DB.QueryRowContext(ctx, queryToGetDriveInfo, driveID)
	var drive models.Drive
	err := row.Scan(&drive.CompanyID, &drive.Company.Name, &drive.Company.Overview, &drive.Company.ContactEmail, &drive.Company.ContactNumber, &drive.Company.LinkedIn, &drive.Company.Website, &drive.DateOfDrive, &drive.DriveDuration, &drive.Location, &drive.Qualifications, &drive.PointsToNote, &drive.JobDescription)

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

func (db *Database) GetAllCompanies(args ...string) ([]models.Company, error) {
	companies := make([]models.Company, 0)

	offset := "1"
	name := ""
	if len(args) > 0 {
		offset = args[0]
		name = args[1]
	}

	var query string
	queryArgs := []interface{}{}

	if name != "" {
		query += "name LIKE ?"
		queryArgs = append(queryArgs, "%"+name+"%")
	}

	queryArgs = append(queryArgs, offset)

	if query != "" {
		query = "SELECT company_id, name, overview, hr_name, contact_email, contact_number, linked_in, website FROM company WHERE " + query + " LIMIT 10 OFFSET ?;"
	} else {
		query = "SELECT company_id, name, overview, hr_name, contact_email, contact_number, linked_in, website FROM company LIMIT 10 OFFSET ?"
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*2)
	defer cancel()

	rows, err := db.DB.QueryContext(ctx, query, queryArgs...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		company := new(models.Company)
		if err := rows.Scan(&company.CompanyID, &company.Name, &company.Overview, &company.HRName, &company.ContactEmail, &company.ContactNumber, &company.LinkedIn, &company.Website); err != nil {
			return nil, err
		}
		companies = append(companies, *company)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return companies, nil
}

func (db *Database) GetAllCompaniesForUser(args ...string) ([]models.CompanyResponse, error) {
	companies := make([]models.CompanyResponse, 0)

	offset := "1"
	name := ""
	if len(args) > 0 {
		offset = args[0]
		name = args[1]
	}

	var query string
	queryArgs := []interface{}{}

	if name != "" {
		query += "name LIKE ?"
		queryArgs = append(queryArgs, "%"+name+"%")
	}

	queryArgs = append(queryArgs, offset)

	if query != "" {
		query = "SELECT * FROM company WHERE " + query + " LIMIT 10 OFFSET ?;"
	} else {
		query = "SELECT * FROM company LIMIT 10 OFFSET ?"
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*2)
	defer cancel()

	rows, err := db.DB.QueryContext(ctx, query, queryArgs...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		company := new(models.CompanyResponse)
		if err := rows.Scan(&company.CompanyID, &company.Name, &company.Overview, &company.LinkedIn, &company.Website); err != nil {
			return nil, err
		}
		companies = append(companies, *company)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return companies, nil
}

func (db *Database) GetCompanyFromCompnayID(companyID string) (interface{}, error) {

	queryToGetCompany := `
    SELECT name, overview, linked_in, website
    FROM company
    WHERE company_id = ?;
    `

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*2)
	defer cancel()

	var company models.CompanyResponse
	row := db.DB.QueryRowContext(ctx, queryToGetCompany, companyID)
	err := row.Scan(&company.Name, &company.Overview, &company.LinkedIn, &company.Website)
	if err != nil {
		return nil, err
	}
	return company, err
}
