package pkg

import (
	"fmt"
	"time"
)

func ParseDates(driveDateStr, deadlineStr string) (time.Time, time.Time, error) {
	// Parse the date strings into time.Time objects
	driveDate, err := time.Parse(time.RFC3339, driveDateStr)
	if err != nil {
		fmt.Println("Error parsing drive_date:", err)
		return time.Time{}, time.Time{}, err
	}
	deadline, err := time.Parse(time.RFC3339, deadlineStr)
	if err != nil {
		fmt.Println("Error parsing deadline:", err)
		return time.Time{}, time.Time{}, err
	}

	driveDateUTC := driveDate.UTC()
	deadlineUTC := deadline.UTC()

	return driveDateUTC, deadlineUTC, nil
}

func ConvertToIST(utcTime time.Time) time.Time {
	istLocation, err := time.LoadLocation("Asia/Kolkata")
	if err != nil {
		fmt.Println("Error loading location:", err)
		return utcTime
	}
	return utcTime.In(istLocation)
}
