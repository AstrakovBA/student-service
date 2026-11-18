package student

import "errors"

var (
	ErrStudentNotFound = errors.New("student not found")
	ErrCmdRepository   = errors.New("unable to command repository")
	ErrQueryRepository = errors.New("unable to query repository")
)
