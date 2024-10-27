package responses

const (
	IvalidJobPosting   HttpMessage = "Failed to post opportunity."
	JobPostingFailed   HttpMessage = "Failed to post opportunity."
	RolesInsertionFail HttpMessage = "Failed to insert role."
)

const (
	DriveNotFound HttpMessage = "No Drive found for the given ID."
	DriveFound    HttpMessage = "Drive available given ID."

	DriveCreated HttpMessage = "Drive created successfully."
)
