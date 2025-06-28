package pkg

import (
	"slices"
	"strings"

	"github.com/DevSoc-exe/placement-portal-backend/internal/models"
)

func DriveDTOMapper(driveBody models.DriveBody) models.Drive {
	drive := models.Drive{
		CompanyID:      driveBody.CompanyID,
		DateOfDrive:    driveBody.DateOfDrive,
		DriveDuration:  driveBody.DriveDuration,
		Roles:          driveBody.Roles,
		MinCGPA:        driveBody.MinCGPA,
		Deadline:       driveBody.Deadline,
		Location:       driveBody.Location,
		Qualifications: driveBody.Qualifications,
		PointsToNote:   driveBody.PointsToNote,
		JobDescription: driveBody.JobDescription,
		DriveType:      driveBody.DriveType,
		RequiredData:   driveBody.RequiredData,
	}

	allowedBranches := strings.Split(driveBody.AllowedBranches, ",")
	drive.Cse_allowed = slices.Contains(allowedBranches, "Computer Science and Engineering")
	drive.Ece_allowed = slices.Contains(allowedBranches, "Electronics and Communication Engineering")
	drive.Mech_allowed = slices.Contains(allowedBranches, "Mechanical Engineering")
	drive.Civ_allowed = slices.Contains(allowedBranches, "Civil Engineering")

	return drive
}

