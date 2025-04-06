package database

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/DevSoc-exe/placement-portal-backend/internal/models"
	_ "github.com/go-sql-driver/mysql"
)

func (s *Database) createUserTable() error {
	query := `CREATE TABLE IF NOT EXISTS users (
    id VARCHAR(36) PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
	phone_number VARCHAR(15) NOT NULL,
	email VARCHAR(255) NOT NULL UNIQUE,
	gender ENUM('MALE', 'FEMALE', 'OTHERS') NOT NULL,
	otp TEXT,
    branch VARCHAR(100) NOT NULL,
    rollnum VARCHAR(100) NOT NULL UNIQUE,
    year_of_admission INT NOT NULL,
    is_verified BOOLEAN NOT NULL DEFAULT FALSE,
    verification_token VARCHAR(255),
    student_type ENUM('Regular', 'PU MEET', 'LEET') NOT NULL,
    created_at DATETIME,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    refresh_token TEXT,
    role ENUM('STUDENT', 'ADMIN', 'MODERATOR') NOT NULL DEFAULT 'STUDENT',
	isOnboarded BOOLEAN NOT NULL DEFAULT FALSE
);
`
	_, err := s.DB.Exec(query)
	return err
}

func (db *Database) CreateUser(user *models.User) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*2)
	defer cancel()

	query := `
    INSERT INTO users (id, name, phone_number, email, rollnum, year_of_admission, branch, student_type, verification_token, role, gender, created_at, updated_at)
    VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, NOW(), NOW());
    `

	_, err := db.DB.ExecContext(ctx, query, user.ID, user.Name, user.PhoneNumber, user.Email, user.RollNumber, user.YearOfAdmission, user.Branch, user.StudentType, user.VerificationToken.String, user.Role, user.Gender)
	if err != nil {
		return err
	}

	return nil
}

func (db *Database) GetUserByEmail(email string) (*models.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*2)
	defer cancel()

	query := `SELECT id, name, phone_number, email, otp, gender, rollnum, year_of_admission, branch, student_type, is_verified, verification_token, role, isOnboarded FROM users WHERE email = ?;`
	row := db.DB.QueryRowContext(ctx, query, email)

	var user models.User
	err := row.Scan(
		&user.ID,
		&user.Name,
		&user.PhoneNumber,
		&user.Email,
		&user.Otp,
		&user.RollNumber,
		&user.Gender,
		&user.YearOfAdmission,
		&user.Branch,
		&user.StudentType,
		&user.IsVerified,
		&user.VerificationToken,
		&user.Role,
		&user.IsOnboarded,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("no user found with email: %s", email)
		}
		return nil, err
	}

	return &user, nil
}

func (db *Database) GetUserByID(id string) (*models.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*2)
	defer cancel()

	query := `SELECT id, name, phone_number, email, otp, gender, rollnum, year_of_admission, branch, student_type, is_verified, verification_token, role, refresh_token, isOnboarded FROM users WHERE id = ?;`
	row := db.DB.QueryRowContext(ctx, query, id)

	var user models.User
	err := row.Scan(
		&user.ID,
		&user.Name,
		&user.PhoneNumber,
		&user.Email,
		&user.Otp,
		&user.Gender,
		&user.RollNumber,
		&user.YearOfAdmission,
		&user.Branch,
		&user.StudentType,
		&user.IsVerified,
		&user.VerificationToken,
		&user.Role,
		&user.RefreshToken,
		&user.IsOnboarded,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("no user found with id: %s", id)
		}
		return nil, err
	}

	return &user, nil
}

func (db *Database) RevokeUserRefreshToken(id string) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*2)
	defer cancel()

	query := `UPDATE users SET refresh_token = ? WHERE id = ?`
	_, err := db.DB.ExecContext(ctx, query, "", id)

	return err
}

func (db *Database) UpdateUserRefreshToken(refreshToken, userId string) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*2)
	defer cancel()

	query := `
		UPDATE users
		SET refresh_token = ?
		WHERE id = ?;
	`

	result, err := db.DB.ExecContext(ctx, query, refreshToken, userId)
	if err != nil {
		return fmt.Errorf("failed to update token: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to check affected rows: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("no user found for the id")
	}

	return nil
}

func (db *Database) SaveOTP(otpToken string, userId string) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*2)
	defer cancel()

	query := `
		UPDATE users
		SET otp = ?
		WHERE id = ?;
	`

	result, err := db.DB.ExecContext(ctx, query, otpToken, userId)
	if err != nil {
		return fmt.Errorf("failed to update password: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to check affected rows: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("no user found for the id")
	}

	return nil
}

func (db *Database) ClearOTP(userID string) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*2)
	defer cancel()

	query := `
		UPDATE users
		SET otp = NULL
		WHERE id = ?;
	`

	result, err := db.DB.ExecContext(ctx, query, userID)
	if err != nil {
		return fmt.Errorf("failed to update password: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to check affected rows: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("no user found for the id")
	}

	return nil
}

func (db *Database) VerifyUser(userId, token string) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*2)
	defer cancel()

	query := `
		UPDATE users
		SET is_verified = 1, verification_token = NULL
		WHERE id = ? AND verification_token = ? AND is_verified = 0;
	`

	result, err := db.DB.ExecContext(ctx, query, userId, token)
	if err != nil {
		return fmt.Errorf("failed to verify user: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to check affected rows: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("invalid token or user already verified")
	}

	return nil
}

func (db *Database) GetAllStudents(args ...string) ([]*models.UserResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*2)
	defer cancel()

	offset := "1"
	gender := ""
	branch := ""
	role := ""
	name := ""
	if len(args) > 0 {
		offset = args[0]
		gender = args[1]
		branch = args[2]
		role = args[3]
		name = args[4]
	}

	var query string
	queryArgs := []interface{}{}

	// Build gender filter
	if gender != "" {
		genderList := strings.Split(gender, ",")
		genderPlaceholders := strings.Repeat("?,", len(genderList)-1) + "?"
		queryArgs = append(queryArgs, convertToInterfaceSlice(genderList)...)
		query += fmt.Sprintf("gender IN (%s)", genderPlaceholders)
	}

	// Build branch filter
	if branch != "" {
		if query != "" {
			query += " AND "
		}
		branchList := strings.Split(branch, ",")
		branchPlaceholders := strings.Repeat("?,", len(branchList)-1) + "?"
		queryArgs = append(queryArgs, convertToInterfaceSlice(branchList)...)
		query += fmt.Sprintf("branch IN (%s)", branchPlaceholders)
	}

	// Build branch filter
	if role != "" {
		if query != "" {
			query += " AND "
		}
		roleList := strings.Split(role, ",")
		rolePlaceholders := strings.Repeat("?,", len(roleList)-1) + "?"
		queryArgs = append(queryArgs, convertToInterfaceSlice(roleList)...)
		query += fmt.Sprintf("role IN (%s)", rolePlaceholders)
	}

	// Build branch filter
	if name != "" {
		if query != "" {
			query += " AND "
		}
		query += "name LIKE ?"
		queryArgs = append(queryArgs, "%"+name+"%")
	}

	queryArgs = append(queryArgs, offset)

	// Base query with filters and pagination
	if query != "" {
		query = "SELECT id, branch, email, gender, name, rollnum, student_type, year_of_admission, isOnboarded, phone_number, role FROM users WHERE " + query + " LIMIT 10 OFFSET ?"
	} else {
		query = "SELECT id, branch, email, gender, name, rollnum, student_type, year_of_admission, isOnboarded, phone_number, role FROM users LIMIT 10 OFFSET ?"
	}

	rows, err := db.DB.QueryContext(ctx, query, queryArgs...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var studentData []*models.UserResponse
	for rows.Next() {
		var data models.UserResponse
		err := rows.Scan(
			&data.ID,
			&data.Branch,
			&data.Email,
			&data.Gender,
			&data.Name,
			&data.RollNumber,
			&data.StudentType,
			&data.YearOfAdmission,
			&data.IsOnboarded,
			&data.PhoneNumber,
			&data.Role,
		)
		if err != nil {
			return nil, err
		}
		studentData = append(studentData, &data)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return studentData, nil
}

func (db *Database) ToggleUserRole(id string) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*2)
	defer cancel()

	query := `UPDATE users SET role = CASE WHEN role = 'STUDENT' THEN 'MODERATOR' ELSE 'STUDENT' END WHERE id = ?;`
	_, err := db.DB.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}

	return nil
}

// Helper function to convert string slice to interface slice
func convertToInterfaceSlice(strs []string) []interface{} {
	interfaces := make([]interface{}, len(strs))
	for i, s := range strs {
		interfaces[i] = s
	}
	return interfaces
}

func (db *Database) GetUserMailsByBranchesAboveCGPA(branches []string, cgpaLimit float32) ([]string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*2)
	defer cancel()
	placeholders := make([]string, len(branches))
	for i := range branches {
		placeholders[i] = "?"
	}

	query := fmt.Sprintf(
		"SELECT users.email FROM users JOIN student_data ON users.id = student_data.id WHERE branch IN (%s) AND cgpa >= %f;",
		strings.Join(placeholders, ","),
		cgpaLimit,
	)

	fmt.Println("Printing the Query")
	fmt.Println(query)

	rows, err := db.DB.QueryContext(ctx, query, convertToInterfaceSlice(branches)...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var mailingList []string
	for rows.Next() {
		var studentEmail string
		err := rows.Scan(&studentEmail)
		if err != nil {
			return nil, err
		}
		mailingList = append(mailingList, studentEmail)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return mailingList, nil
}
