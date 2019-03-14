package skelerror

type ServiceError struct {
	key string //length which caused the error
	err string //error description
}

func (e *ServiceError) Error() string {
	return e.err
}

func (e *ServiceError) Key() string {
	return e.key
}

var (
	// ErrBadRouting is returned when an expected path variable is missing.
	InvalidCredentials = ServiceError{"INVALID_CREDENTIALS", "Invalid Credentials provided"}
	JSONDecodeError    = ServiceError{"JSON_DECODE_ERROR", "Unable to decode JSON"}
	NoBody             = ServiceError{"NO_BODY", "No request body found"}
	QueryFailure       = ServiceError{"QUERY_FAILURE", "Error querying DB"}
)
