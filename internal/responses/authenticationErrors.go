package responses

const (
	Unauthorized              HttpMessage = "Access denied. Please log in to continue."
	FailedRefreshToken        HttpMessage = "Unable to refresh your session. Please log in again."
	ErrorTokenCreation        HttpMessage = "There was an issue creating your token. Please try again."
	ErrorRefreshTokenCreation HttpMessage = "There was an issue creating your refresh token. Please try again."
	LoginSuccess              HttpMessage = "Welcome back! You have successfully logged in."
)
