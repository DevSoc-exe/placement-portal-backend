package database

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/DevSoc-exe/placement-portal-backend/internal/models"
	"github.com/DevSoc-exe/placement-portal-backend/internal/models/dto"
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
		drive_date DATETIME NOT NULL,
		drive_duration INT NOT NULL,
		location VARCHAR(255),
		min_cgpa DECIMAL(3,2) NOT NULL,
		deadline DATETIME NOT NULL,
		qualifications LONGTEXT NOT NULL,
		points_to_note LONGTEXT NOT NULL,
		job_description VARCHAR(255),
		drive_type enum('on-campus', 'company-office', 'online') DEFAULT 'on-campus' NOT NULL,
		cse_allowed BOOLEAN DEFAULT FALSE,
		ece_allowed BOOLEAN DEFAULT FALSE,
		civ_allowed BOOLEAN DEFAULT FALSE,
		mech_allowed BOOLEAN DEFAULT FALSE,
		required_data VARCHAR(255) DEFAULT 'name,phone_number,email,branch,rollnum,cgpa,has_backlogs' NOT NULL,
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

func (db *Database) CreateNewDriveUsingObject(driveData models.Drive) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	tx, err := db.DB.BeginTx(ctx, nil)
	if err != nil {
		return "", err
	}

	queryToInsertDrive := `
	INSERT INTO drive (id, company_id, drive_date, drive_duration, location, qualifications, points_to_note, job_description, min_cgpa, deadline, drive_type, required_data, cse_allowed, ece_allowed, civ_allowed, mech_allowed)
	VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?);
	`

	driveUUID := nanoid.New()

	_, err = tx.ExecContext(ctx, queryToInsertDrive, driveUUID, driveData.CompanyID, driveData.DateOfDrive, driveData.DriveDuration, driveData.Location, driveData.Qualifications, driveData.PointsToNote, driveData.JobDescription, driveData.MinCGPA, driveData.Deadline, driveData.DriveType, driveData.RequiredData, driveData.Cse_allowed, driveData.Ece_allowed, driveData.Civ_allowed, driveData.Mech_allowed)
	if err != nil {
		tx.Rollback()
		return "", err
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
		return "", err
	}

	err = tx.Commit()
	if err != nil {
		return "", err
	}

	return driveUUID, err
}

func (db *Database) DeleteJobUsingDriveID(driveID string) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
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
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
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

func (db *Database) GetJobPostingUsingDriveID(driveID string) (*models.Drive, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	queryToGetDriveInfo := `
    SELECT d.id, c.company_id, c.name, c.overview, c.contact_email, c.contact_number,
           c.linked_in, c.website, d.drive_date, d.drive_duration, d.location,
           d.qualifications, d.points_to_note, d.job_description, d.min_cgpa,
           d.deadline, d.drive_type, d.cse_allowed, d.ece_allowed, d.civ_allowed,
           d.mech_allowed, d.required_data
    FROM company c
    JOIN drive d ON c.company_id = d.company_id
    WHERE d.id = ?;
    `

	row := db.DB.QueryRowContext(ctx, queryToGetDriveInfo, driveID)
	fmt.Println(row)

	var drive models.Drive

	err := row.Scan(&drive.ID, &drive.CompanyID, &drive.Company.Name, &drive.Company.Overview,
		&drive.Company.ContactEmail, &drive.Company.ContactNumber, &drive.Company.LinkedIn,
		&drive.Company.Website, &drive.DateOfDrive, &drive.DriveDuration, &drive.Location,
		&drive.Qualifications, &drive.PointsToNote, &drive.JobDescription, &drive.MinCGPA,
		&drive.Deadline, &drive.DriveType, &drive.Cse_allowed, &drive.Ece_allowed,
		&drive.Civ_allowed, &drive.Mech_allowed, &drive.RequiredData)

	if err != nil {
		return nil, err
	}

	// parsedDeadline, err := time.Parse("2006-01-02 15:04:05", deadline)
	// if err != nil {
	// 	return nil, fmt.Errorf("failed to parse deadline: %v", err)
	// }
	// drive.Deadline = parsedDeadline
	// Convert dates from UTC to IST
	// drive.DateOfDrive = pkg.ConvertToIST(drive.DateOfDrive)
	// drive.Deadline = pkg.ConvertToIST(drive.Deadline)

	//! Dont Try to understand this, it's a hack, not my proudest moment
	date := drive.Deadline.UTC().String()

	date = date[0:20] + "+0530 IST"

	parsedDeadline, err := time.Parse("2006-01-02 15:04:05 -0700 MST", date)

	drive.Expired = parsedDeadline.Before(time.Now())

	roles, err := db.GetRolesUsingDriveID(driveID)
	if err != nil {
		drive.Roles = []models.Role{} // Ensure Roles is not nil
	} else {
		drive.Roles = roles
	}
	return &drive, nil
}

func (db *Database) GetAppliedRole(userID string, driveID string) (*models.Role, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	query := `
	SELECT role.id, role.drive_id, title, stipend_low, stipend_high, salary_low, salary_high
	FROM role
	JOIN applications ON role.id = applications.role_id
	WHERE applications.user_id = ? AND role.drive_id = ?;
	`

	row := db.DB.QueryRowContext(ctx, query, userID, driveID)
	role := &models.Role{}
	err := row.Scan(&role.ID, &role.DriveID, &role.Title, &role.StipendLow, &role.StipendHigh, &role.SalaryLow, &role.SalaryHigh)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // No role found, not an error
		}
		return nil, err
	}
	return role, nil
}

func (db *Database) GetAllDrivesForUser() ([]models.DriveResponse, error) {
	drives := make([]models.DriveResponse, 0)

	query := `
		SELECT drive.id, company.name, drive_date, drive_duration, location, qualifications, points_to_note, job_description, min_cgpa, deadline, drive_type
		FROM drive
		JOIN company ON drive.company_id = company.company_id
		ORDER BY drive_date DESC;
	`
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	rows, err := db.DB.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		drive := new(models.DriveResponse)

		if err := rows.Scan(
			&drive.ID,
			&drive.CompanyName,
			&drive.DateOfDrive,
			&drive.DriveDuration,
			&drive.Location,
			&drive.Qualifications,
			&drive.PointsToNote,
			&drive.JobDescription,
			&drive.MinCGPA,
			&drive.Deadline,
			&drive.DriveType,
		); err != nil {
			return nil, err
		}

		// // Parse the deadline string to time.Time
		// parsedDeadline, err := time.Parse("2006-01-02 15:04:05", deadlineStr)
		// if err != nil {
		// 	return nil, err
		// }
		// drive.Deadline = parsedDeadline

		// Retrieve roles as before
		roles, err := db.GetRolesUsingDriveID(drive.ID)
		if err != nil {
			return nil, err
		}
		drive.Roles = roles

		drives = append(drives, *drive)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return drives, nil
}

func (db *Database) GetRolesUsingDriveID(driveID string) ([]models.Role, error) {
	roles := make([]models.Role, 0)

	queryToGetRoles := `
    SELECT id, drive_id, title, stipend_low, stipend_high, salary_low, salary_high
    FROM role
    WHERE drive_id = ?;
    `

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
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
		query = "SELECT company_id, name, overview, linked_in, website FROM company WHERE " + query + " LIMIT 10 OFFSET ?;"
	} else {
		query = "SELECT company_id, name, overview, linked_in, website FROM company LIMIT 10 OFFSET ?"
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

func (db *Database) GetDriveApplicantsForRole(roleID, required_data, driveID string) (*sql.Rows, []string, error) {

	// Split the requested columns
	columnList := strings.Split(required_data, ",")

	// Construct the query dynamically
	query := fmt.Sprintf(`
		SELECT %s
		FROM applications
		JOIN users ON applications.user_id = users.id
		JOIN student_data ON users.id = student_data.id
		WHERE applications.role_id = ? AND applications.drive_id = ?`, strings.Join(columnList, ", "))

	// Execute the query
	rows, err := db.DB.Query(query, roleID, driveID)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to execute query: %w", err)
	}

	// Get column names dynamically
	columns, err := rows.Columns()
	if err != nil {
		return nil, nil, fmt.Errorf("failed to get column names: %w", err)
	}

	return rows, columns, nil
}

func (db *Database) GetCompanyUsingCompanyID(companyID string) (*models.CompanyResponse, error) {
	query := "SELECT company_id, name, overview, linked_in, website FROM company WHERE company_id = ?"
	queryArgs := []interface{}{companyID}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*2)
	defer cancel()

	row := db.DB.QueryRowContext(ctx, query, queryArgs...)
	company := new(models.CompanyResponse)

	if err := row.Scan(&company.CompanyID, &company.Name, &company.Overview, &company.LinkedIn, &company.Website); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return company, nil
}

func (db *Database) GetApplicantsForDrive(driveID string) ([]dto.StudentApplicationDTO, error) {

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	query := `
		SELECT u.id, a.id, u.name, u.branch , u.rollnum,  u.email , u.gender , r.title , a.applied_at, r.id, a.is_placed
		FROM users AS u
		JOIN applications AS a
		ON u.id = a.user_id
		JOIN role r
		ON a.role_id = r.id
		WHERE a.drive_id = ? and a.deleted is null;
	`

	rows, err := db.DB.QueryContext(ctx, query, driveID)
	if err != nil {
		return nil, fmt.Errorf("failed to execute query: %w", err)
	}

	defer rows.Close()

	students := make([]dto.StudentApplicationDTO, 0)

	for rows.Next() {
		var student dto.StudentApplicationDTO

		if err := rows.Scan(
			&student.Id,
			&student.ApplicationId,
			&student.Name,
			&student.Branch,
			&student.RollNum,
			&student.Email,
			&student.Gender,
			&student.Role,
			&student.AppliedAt,
			&student.RoleId,
			&student.IsPlaced,
		); err != nil {
			return nil, err
		}

		students = append(students, student)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return students, nil
}
