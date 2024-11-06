package database

import (
	"context"
	"fmt"
	"time"

	"github.com/DevSoc-exe/placement-portal-backend/internal/models"
)

func (s *Database) createStudentDataTable() error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*2)
	defer cancel()

	query := `
		CREATE TABLE IF NOT EXISTS student_data (
		id VARCHAR(36) PRIMARY KEY NOT NULL,
		sgpasem1 DECIMAL(3, 2) NOT NULL,
		sgpasem2 DECIMAL(3, 2) NOT NULL,
		sgpasem3 DECIMAL(3, 2) NOT NULL,
		sgpasem4 DECIMAL(3, 2) NOT NULL,
		sgpasem5 DECIMAL(3, 2) NOT NULL,
		sgpasem6 DECIMAL(3, 2) NOT NULL,
		cgpa DECIMAL(3, 2) NOT NULL,
		marks10th DECIMAL(5, 2) NOT NULL,
		marks12th DECIMAL(5, 2) NOT NULL,
		sgpa_proofs VARCHAR(255) NOT NULL,
		achievement_certificates VARCHAR(255) NOT NULL,
		college_id_card VARCHAR(255) NOT NULL,
		CONSTRAINT fk_user FOREIGN KEY (id) REFERENCES users(id)
	);
	`

	_, err := s.DB.ExecContext(ctx, query)
	if err != nil {
		return fmt.Errorf("failed to create student data table: %w", err)
	}

	return nil
}

func (db *Database) AddStudentData(user *models.StudentData) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	tx, err := db.DB.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	query := `insert into student_data (id, sgpasem1, sgpasem2, sgpasem3, sgpasem4, sgpasem5, sgpasem6,cgpa, marks10th, marks12th, sgpa_proofs, achievement_certificates, college_id_card) values (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?);`
	_, err = tx.ExecContext(ctx, query, user.ID, user.Sem1SGPA, user.Sem2SGPA, user.Sem3SGPA, user.Sem4SGPA, user.Sem5SGPA, user.Sem6SGPA, user.Cgpa, user.Marks10th, user.Marks12th, user.SgpaProofs, user.AchievementCertificates, user.CollegeIdCard)
	if err != nil {
		tx.Rollback()
		return err
	}

	query = `update users set isOnboarded = 1 where id = ?;`
	_, err = tx.ExecContext(ctx, query, user.ID)
	if err != nil {
		tx.Rollback()
		return err
	}

	err = tx.Commit()
	if err != nil {
		return err
	}

	return nil
}

func (db *Database) GetStudentDataByID(id string) (*models.StudentData, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*2)
	defer cancel()

	query := `select * from student_data where id = ?;`
	row := db.DB.QueryRowContext(ctx, query, id)

	var studentData models.StudentData

	err := row.Scan(&studentData.ID, &studentData.Sem1SGPA, &studentData.Sem2SGPA, &studentData.Sem3SGPA, &studentData.Sem4SGPA, &studentData.Sem5SGPA, &studentData.Sem6SGPA, &studentData.Cgpa, &studentData.Marks10th, &studentData.Marks12th, &studentData.SgpaProofs, &studentData.AchievementCertificates, &studentData.CollegeIdCard)
	if err != nil {
		return nil, err
	}

	return &studentData, nil
}

func (db *Database) UpdateStudentData(user *models.StudentData) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*2)
	defer cancel()

	query := `update student_data set sgpasem1 = ?, sgpasem2 = ?, sgpasem3 = ?, sgpasem4 = ?, sgpasem5 = ?, sgpasem6 = ?, cgpa = ?, marks10th = ?, marks12th = ?, sgpa_proofs = ?, achievement_certificates = ?, college_id_card = ? where id = ?;`

	_, err := db.DB.ExecContext(ctx, query, user.Sem1SGPA, user.Sem2SGPA, user.Sem3SGPA, user.Sem4SGPA, user.Sem5SGPA, user.Sem6SGPA, user.Cgpa, user.Marks10th, user.Marks12th, user.SgpaProofs, user.AchievementCertificates, user.CollegeIdCard, user.ID)

	if err != nil {
		return err
	}

	return nil
}

func (db *Database) DeleteStudentData(id string) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*2)
	defer cancel()

	tx, err := db.DB.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	query := `delete from student_data where id = ?;`
	_, err = tx.ExecContext(ctx, query, id)
	if err != nil {
		tx.Rollback()
		return err
	}

	query = `update users set isOnboarded = 0 where id = ?;`
	_, err = tx.ExecContext(ctx, query, id)
	if err != nil {
		tx.Rollback()
		return err
	}

	err = tx.Commit()
	if err != nil {
		return err
	}

	return nil
}

func (db *Database) GetAllStudentData(args ...string) ([]*models.StudentData, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*2)
	defer cancel()

	offset := "1"
	if len(args) > 0 {
		offset = args[0]
	}

	query := `SELECT * FROM student_data LIMIT 10 OFFSET $1;`
	rows, err := db.DB.QueryContext(ctx, query, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var studentData []*models.StudentData

	for rows.Next() {
		var data models.StudentData
		err := rows.Scan(&data.ID, &data.Sem1SGPA, &data.Sem2SGPA, &data.Sem3SGPA, &data.Sem4SGPA, &data.Sem5SGPA, &data.Sem6SGPA, &data.Cgpa, &data.Marks10th, &data.Marks12th, &data.SgpaProofs, &data.AchievementCertificates, &data.CollegeIdCard)
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
