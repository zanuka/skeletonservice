package skelerror

type AuthServiceError struct {
	key string //length which caused the error
	err string //error description
}

func (e *AuthServiceError) Error() string {
	return e.err
}

func (e *AuthServiceError) Key() string {
	return e.key
}

var (
	// ErrBadRouting is returned when an expected path variable is missing.
	InvalidCredentials = AuthServiceError{"INVALID_CREDENTIALS", "Invalid Credentials provided"}
	JSONDecodeError    = AuthServiceError{"JSON_DECODE_ERROR", "Unable to decode JSON"}
	NoBody             = AuthServiceError{"NO_BODY", "No request body found"}
	QueryFailure       = AuthServiceError{"QUERY_FAILURE", "Error querying DB"}
)
