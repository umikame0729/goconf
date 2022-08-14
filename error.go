package goconf

import "fmt"

type ErrorIsUpdated struct {
	Path    string
	BakPath string
}

func (e *ErrorIsUpdated) Error() string {
	return fmt.Sprintf("`%s` updated, backup at `%s`", e.Path, e.BakPath)
}

func IsUpdated(err error) bool {
	_, ok := err.(*ErrorIsUpdated)
	return ok
}

type ErrorIsNewCreated struct {
	Err error
}

func (e *ErrorIsNewCreated) Error() string {
	return e.Err.Error()
}

func IsNewCreated(err error) bool {
	_, ok := err.(*ErrorIsNewCreated)
	return ok
}
