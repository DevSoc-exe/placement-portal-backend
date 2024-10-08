package responses

type HttpMessage string

const (
	TryAgain            HttpMessage = "Something went wrong. Please try again!"
	InternalServerError HttpMessage = "Oops! Something went wrong on our end. Please try again later."
	DatabaseError       HttpMessage = "There was an issue connecting to the database. Please try again later."
)
