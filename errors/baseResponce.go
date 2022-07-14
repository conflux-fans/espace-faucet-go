package errors

type BaseError struct {
	Code    int
	Message string
}

func NewError(code int, message string) BaseError {
	return BaseError{
		Code:    code,
		Message: message,
	}
}

func (err BaseError) Error() string {
	//return fmt.Sprintf("%d %s", err.Code, err.Message)
	return err.Message
}
