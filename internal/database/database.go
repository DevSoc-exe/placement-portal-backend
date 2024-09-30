package database

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/DevSoc-exe/placement-portal-backend/internal/models"
	_ "github.com/go-sql-driver/mysql"
)

type Database struct {
	DB *sql.DB
}

func CreateDatabase(db *sql.DB) *Database {
	return &Database{
		DB: db,
	}
}

func ConnectToDB(dsn string) (*sql.DB, error) {
	// Open a connection to the database
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}

	// Verify the connection to the database
	if err := db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}

func (s *Database) InitDB() error {
	err := s.createUserTable()
	if err != nil {
		return fmt.Errorf("could not create tables: %s", err)
	}

	return nil
}

func (s *Database) createUserTable() error {
	query := `CREATE TABLE IF NOT EXISTS users (
    id VARCHAR(36) PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    email VARCHAR(255) NOT NULL UNIQUE,
    password VARCHAR(255) NOT NULL,
    branch VARCHAR(100) NOT NULL,
    rollnum VARCHAR(100) NOT NULL UNIQUE,
    year_of_admission INT NOT NULL,
    is_verified BOOLEAN NOT NULL DEFAULT FALSE,
    verification_token VARCHAR(255),
    student_type ENUM('Regular', 'PU MEET', 'LEET') NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    refresh_token VARCHAR(255),
    role ENUM('STUDENT', 'ADMIN', 'MODERATOR') NOT NULL DEFAULT 'STUDENT'
);
`

	_, err := s.DB.Exec(query)
	return err
}

func (db *Database) CreateUser(user *models.User) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*2)
	defer cancel()

	query := `
    INSERT INTO users (id, name, email, password, rollnum, year_of_admission, branch, student_type, verification_token, role, created_at, updated_at)
    VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, NOW(), NOW());
    `

	_, err := db.DB.ExecContext(ctx, query, user.ID, user.Name, user.Email, user.Password, user.RollNumber, user.YearOfAdmission, user.Branch, user.StudentType, user.VerificationToken, user.Role)
	if err != nil {
		return err
	}

	return nil
}

func (db *Database) GetUserByEmail(email string) (*models.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*2)
	defer cancel()

	query := `SELECT id, name, email, password, rollnum, year_of_admission, branch, student_type, is_verified, verification_token, role FROM users WHERE email = ?;`
	row := db.DB.QueryRowContext(ctx, query, email)

	var user models.User
	err := row.Scan(&user.ID, &user.Name, &user.Email, &user.Password, &user.RollNumber, &user.YearOfAdmission, &user.Branch, &user.StudentType, &user.IsVerified, &user.VerificationToken, &user.Role)
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

	query := `SELECT id, name, email, password, rollnum, year_of_admission, branch, student_type, is_verified, verification_token, role FROM users WHERE id = ?;`
	row := db.DB.QueryRowContext(ctx, query, id)

	var user models.User
	err := row.Scan(&user.ID, &user.Name, &user.Email, &user.Password, &user.RollNumber, &user.YearOfAdmission, &user.Branch, &user.StudentType, &user.IsVerified, &user.VerificationToken, &user.Role)
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

func (db *Database) UpdateUserRefreshToken(refresh_token, userId string) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*2)
	defer cancel()

	query := `
		UPDATE users
		SET refresh_token = ?,
		WHERE id = ?;
	`
	var verifiedUserId models.User
	err := db.DB.QueryRowContext(ctx, query, refresh_token, userId).Scan(&verifiedUserId)
	if err != nil {
		if err == sql.ErrNoRows {
			return fmt.Errorf("Some Error Occured!")
		}
		return fmt.Errorf("failed to verify user: %w", err)
	}

	return nil
}

func (db *Database) VerifyUser(userId, token string) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*2)
	defer cancel()

	query := `
		UPDATE users
		SET is_verified = TRUE, verification_token = NULL
		WHERE id = ? AND verification_token = ? AND is_verified = FALSE;
	`

	var verifiedUserId models.User
	err := db.DB.QueryRowContext(ctx, query, userId, token).Scan(&verifiedUserId)
	if err != nil {
		if err == sql.ErrNoRows {
			return fmt.Errorf("invalid token or user already verified")
		}
		return fmt.Errorf("failed to verify user: %w", err)
	}

	return nil
}
