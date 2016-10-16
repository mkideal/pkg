package optvar

import (
	"errors"
)

var (
	ErrDataSourceIsNil = errors.New("data source is nil")
)

type MissingRequiredVarError string

func (e MissingRequiredVarError) Error() string {
	return "missing required `" + string(e) + "`"
}

func MissingRequiredVar(name string) error {
	return MissingRequiredVarError(name)
}

func getError(required bool, src *Source, name string) error {
	if src != nil {
		return src.Err
	} else if required {
		return MissingRequiredVar(name)
	}
	return nil
}
