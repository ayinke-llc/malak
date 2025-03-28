// Code generated by go-enum DO NOT EDIT.
// Version:
// Revision:
// Build Date:
// Built By:

package malak

import (
	"errors"
	"fmt"
)

const (
	// DashboardLinkTypeDefault is a DashboardLinkType of type default.
	DashboardLinkTypeDefault DashboardLinkType = "default"
	// DashboardLinkTypeContact is a DashboardLinkType of type contact.
	DashboardLinkTypeContact DashboardLinkType = "contact"
)

var ErrInvalidDashboardLinkType = errors.New("not a valid DashboardLinkType")

// String implements the Stringer interface.
func (x DashboardLinkType) String() string {
	return string(x)
}

// IsValid provides a quick way to determine if the typed value is
// part of the allowed enumerated values
func (x DashboardLinkType) IsValid() bool {
	_, err := ParseDashboardLinkType(string(x))
	return err == nil
}

var _DashboardLinkTypeValue = map[string]DashboardLinkType{
	"default": DashboardLinkTypeDefault,
	"contact": DashboardLinkTypeContact,
}

// ParseDashboardLinkType attempts to convert a string to a DashboardLinkType.
func ParseDashboardLinkType(name string) (DashboardLinkType, error) {
	if x, ok := _DashboardLinkTypeValue[name]; ok {
		return x, nil
	}
	return DashboardLinkType(""), fmt.Errorf("%s is %w", name, ErrInvalidDashboardLinkType)
}
