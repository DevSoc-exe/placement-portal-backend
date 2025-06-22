package responses

const (
	InvalidCredentials HttpMessage = "The entered password or email is incorrect. Please Try Again!"
	LoginSuccessful HttpMessage = "Login Successful!"
	UserNotFound      HttpMessage = "We couldn't find a user with the provided information."
	FailedToHash  HttpMessage = "There was an issue securing your password. Please try again."
	UsernameTaken HttpMessage = "This username is already in use. Please choose another."
	EmailTaken    HttpMessage = "An account with this email already exists. Please use a different email."
)
