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
		is_placed INT DEFAULT 0,
		deleted TIMESTAMP DEFAULT NULL,
		last_updated_by VARCHAR(36),
		updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
		applied_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
		FOREIGN KEY (role_id) REFERENCES role(id) ON DELETE CASCADE,
		FOREIGN KEY (drive_id) REFERENCES drive(id) ON DELETE CASCADE,
		CONSTRAINT fk_last_updated_by FOREIGN KEY (last_updated_by) REFERENCES users(id),
		CONSTRAINT chk_is_placed CHECK (is_placed IN (-1, 0, 1)),
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

	loc, err := time.LoadLocation("Asia/Kolkata")
	if err != nil {
		return fmt.Errorf("could not load IST location: %s", err.Error())
	}
	currentTime := time.Now().In(loc)

	query := `
	INSERT INTO applications (id, user_id, role_id, drive_id, applied_at, last_updated_by)
	VALUES (?, ?, ?, ?, ?, ?);
	`
	_, err = db.DB.ExecContext(ctx, query, nanoid.New(), userID, roleID, driveID, currentTime, userID)
	if err != nil {
		return fmt.Errorf("could not apply for drive: %s", err.Error())
	}

	return nil
}

func (db *Database)	MarkStudentAsPlaced(updatedById string, applicationId string) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*2)
	defer cancel()

	query := `UPDATE applications a
			SET a.is_placed = 1, a.last_updated_by = ?, a.updated_at = NOW()
			where a.id = ?;
			`

	_, err := db.DB.ExecContext(ctx, query, updatedById, applicationId)
	if err != nil {
		return err
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
