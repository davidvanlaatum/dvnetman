package auth

type NotLoggedInError struct {
}

func (e *NotLoggedInError) Error() string {
	return "Not logged in"
}

type NotAuthorizedError struct {
}

func (e *NotAuthorizedError) Error() string {
	return "Not authorized"
}
