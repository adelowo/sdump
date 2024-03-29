// Code generated by go-enum DO NOT EDIT.
// Version:
// Revision:
// Build Date:
// Built By:

package config

import (
	"errors"
	"fmt"
)

const (
	// DatabaseTypePsql is a DatabaseType of type psql.
	DatabaseTypePsql DatabaseType = "psql"
)

var ErrInvalidDatabaseType = errors.New("not a valid DatabaseType")

// String implements the Stringer interface.
func (x DatabaseType) String() string {
	return string(x)
}

// IsValid provides a quick way to determine if the typed value is
// part of the allowed enumerated values
func (x DatabaseType) IsValid() bool {
	_, err := ParseDatabaseType(string(x))
	return err == nil
}

var _DatabaseTypeValue = map[string]DatabaseType{
	"psql": DatabaseTypePsql,
}

// ParseDatabaseType attempts to convert a string to a DatabaseType.
func ParseDatabaseType(name string) (DatabaseType, error) {
	if x, ok := _DatabaseTypeValue[name]; ok {
		return x, nil
	}
	return DatabaseType(""), fmt.Errorf("%s is %w", name, ErrInvalidDatabaseType)
}
