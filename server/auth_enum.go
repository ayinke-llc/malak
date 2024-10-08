// Code generated by go-enum DO NOT EDIT.
// Version:
// Revision:
// Build Date:
// Built By:

package server

import (
	"errors"
	"fmt"
)

const (
	// CookieNameUser is a CookieName of type user.
	CookieNameUser CookieName = "user"
)

var ErrInvalidCookieName = errors.New("not a valid CookieName")

// String implements the Stringer interface.
func (x CookieName) String() string {
	return string(x)
}

// IsValid provides a quick way to determine if the typed value is
// part of the allowed enumerated values
func (x CookieName) IsValid() bool {
	_, err := ParseCookieName(string(x))
	return err == nil
}

var _CookieNameValue = map[string]CookieName{
	"user": CookieNameUser,
}

// ParseCookieName attempts to convert a string to a CookieName.
func ParseCookieName(name string) (CookieName, error) {
	if x, ok := _CookieNameValue[name]; ok {
		return x, nil
	}
	return CookieName(""), fmt.Errorf("%s is %w", name, ErrInvalidCookieName)
}
