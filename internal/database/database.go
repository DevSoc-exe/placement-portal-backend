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
		rollnum VARCHAR(100) NOT NULL UNIQUE,
		year_of_admission INT NOT NULL,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
						ON UPDATE CURRENT_TIMESTAMP
	);`

	_, err := s.DB.Exec(query)
	return err
}

func (db *Database) GetUserByEmail(email string) (*models.User, error) {
    ctx, cancel := context.WithTimeout(context.Background(), time.Second*2)
    defer cancel()

    query := `SELECT id, name, email, password, rollnum, year_of_admission FROM users WHERE email = ?;`
    row := db.DB.QueryRowContext(ctx, query, email)

    var user models.User
    err := row.Scan(&user.ID, &user.Name, &user.Email, &user.Password, &user.RollNumber, &user.YearOfAdmission)
    if err != nil {
        if err == sql.ErrNoRows {
            return nil, fmt.Errorf("no user found with email: %s", email)
        }
        return nil, err
    }

    return &user, nil
}



func (db *Database) CreateUser(user *models.User) error {
    ctx, cancel := context.WithTimeout(context.Background(), time.Second*2)
    defer cancel()

    query := `
    INSERT INTO users (id, name, email, password, rollnum, year_of_admission, created_at, updated_at)
    VALUES (?, ?, ?, ?, ?, ?, NOW(), NOW());
    `

    _, err := db.DB.ExecContext(ctx, query,user.ID, user.Name, user.Email, user.Password, user.RollNumber, user.YearOfAdmission)
    if err != nil {
        return err
    }

    return nil
}
