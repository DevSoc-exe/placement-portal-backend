package database

import (
	"context"
	"fmt"
	"time"

	"github.com/aidarkhanov/nanoid"
)

func (db *Database) createApplicationsTable() error {

	// application_status ENUM('PENDING', 'APPROVED', 'REJECTED') DEFAULT 'PENDING',
	query := `CREATE TABLE IF NOT EXISTS applications (
		id VARCHAR(36) PRIMARY KEY,
		user_id VARCHAR(36) NOT NULL,
		role_id VARCHAR(36) NOT NULL,
		drive_id VARCHAR(36) NOT NULL,
		applied_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
		FOREIGN KEY (role_id) REFERENCES role(id) ON DELETE CASCADE,
		FOREIGN KEY (drive_id) REFERENCES drive(id) ON DELETE CASCADE,
		UNIQUE (user_id, drive_id) -- ensures that each user can apply only once per drive
	);
	`
	_, err := db.DB.Exec(query)
	if err != nil {
		return fmt.Errorf("could not create role table: %s", err.Error())
	}

	return nil
}

func (db *Database) ApplyForDrive(userID, roleID, driveID string) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*2)
	defer cancel()

	query := `
	INSERT INTO applications (id, user_id, role_id, drive_id)
	VALUES (?, ?, ?, ?);
	`

	_, err := db.DB.ExecContext(ctx, query, nanoid.New(), userID, roleID, driveID)
	if err != nil {
		return fmt.Errorf("could not apply for drive: %s", err.Error())
	}

	return nil
}

func (db *Database) checkIfUserHasApplied(userID, driveID string) (bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*2)
	defer cancel()

	query := `SELECT id FROM applications WHERE user_id = ? AND drive_id = ?;`

	row := db.DB.QueryRowContext(ctx, query, userID, driveID)
	var id string
	err := row.Scan(&id)
	if err != nil {
		return false, nil
	}

	return true, nil
}
