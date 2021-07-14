package errors

type SystemError struct {
	str        string
	InnerError error
}

func SystemNew(msg string, inner error) error {
	return &SystemError{str: msg, InnerError: inner}
}

func (e *SystemError) Error() string {
	if e.InnerError == nil {
		return e.str
	} else {
		return e.str + ". Inner Error: " + e.InnerError.Error()
	}
}
