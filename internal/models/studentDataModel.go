package models

type StudentData struct {
	ID                      string  `json:"id"`
	Sem1SGPA                float32 `json:"sgpasem1"`
	Sem2SGPA                float32 `json:"sgpasem2"`
	Sem3SGPA                float32 `json:"sgpasem3"`
	Sem4SGPA                float32 `json:"sgpasem4"`
	Sem5SGPA                float32 `json:"sgpasem5"`
	Sem6SGPA                float32 `json:"sgpasem6"`
	Cgpa                    float32 `json:"cgpa"`
	Marks10th               float32 `json:"marks10th"`
	Marks12th               float32 `json:"marks12th"`
	SgpaProofs              string  `json:"sgpaProofs"`
	AchievementCertificates string  `json:"achievementCertificates"`
	CollegeIdCard           string  `json:"collegeIdCard"`
	HasBacklogs			bool    `json:"has_backlogs"`
}
