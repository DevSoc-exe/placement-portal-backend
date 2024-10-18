package handlers

import (
	"context"
	"fmt"

	// "net/http"
	"database/sql"
	// "os"
	"time"

	"github.com/DevSoc-exe/placement-portal-backend/internal/models"
	// "github.com/DevSoc-exe/placement-portal-backend/internal/pkg"
	// "github.com/aidarkhanov/nanoid"
	// "github.com/gin-gonic/gin"
	// "golang.org/x/crypto/bcrypt"
)

type Database struct {
	DB *sql.DB
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
			return nil, fmt.Errorf("no user found with email: %s", id)
		}
		return nil, err
	}
	return &user, nil
}

func (db *Database) UpdateUserByKey(updateContext map[string]interface{}) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*2)
	defer cancel()

	query := `
    UPDATE users 
    SET ` + updateContext["updatedElement"].(string) + ` = ?
    WHERE id = ?;
    `
	result, err := db.DB.ExecContext(ctx, query, updateContext["newValue"], updateContext["id"])
	if err != nil {
		return err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return fmt.Errorf("no user found with id: %v", updateContext["id"])
	}
	return nil
}
