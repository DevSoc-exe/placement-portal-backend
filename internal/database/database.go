package database

import (
	// "context"
	"database/sql"
	"fmt"

	"time"

	// "github.com/DevSoc-exe/placement-portal-backend/internal/models"
	// "github.com/gin-gonic/gin"
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

    db.SetMaxOpenConns(100) 
    db.SetMaxIdleConns(10) 
    db.SetConnMaxLifetime(30 * time.Minute) 

    if err := db.Ping(); err != nil {
        return nil, err
    }

    return db, nil 
}

func (s *Database) InitDB() error {
	err := s.createUserTable()
	if err != nil {
		return fmt.Errorf("could not create users table: %s", err)
	}

	err = s.createJobsTable()
	if err != nil {
		return fmt.Errorf("could not create Jobs table: %s", err)
	}

	err = s.createStudentDataTable()
	if err != nil {
		return fmt.Errorf("could not create student data table: %s", err)
	}

	err = s.createApplicationsTable()
	if err != nil {
		return fmt.Errorf("could not create applications table: %s", err)
	}

	// err := s.createUserTable()
	// if err != nil {
	// 	return fmt.Errorf("could not create users table: %s", err)
	// }

	// err := s.createUserTable()
	// if err != nil {
	// 	return fmt.Errorf("could not create users table: %s", err)
	// }

	return nil
}
