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
	// StatusSuccess is a Status of type Success.
	StatusSuccess Status = iota
	// StatusFailed is a Status of type Failed.
	StatusFailed
)

var ErrInvalidStatus = errors.New("not a valid Status")

const _StatusName = "successfailed"

var _StatusMap = map[Status]string{
	StatusSuccess: _StatusName[0:7],
	StatusFailed:  _StatusName[7:13],
}

// String implements the Stringer interface.
func (x Status) String() string {
	if str, ok := _StatusMap[x]; ok {
		return str
	}
	return fmt.Sprintf("Status(%d)", x)
}

// IsValid provides a quick way to determine if the typed value is
// part of the allowed enumerated values
func (x Status) IsValid() bool {
	_, ok := _StatusMap[x]
	return ok
}

var _StatusValue = map[string]Status{
	_StatusName[0:7]:  StatusSuccess,
	_StatusName[7:13]: StatusFailed,
}

// ParseStatus attempts to convert a string to a Status.
func ParseStatus(name string) (Status, error) {
	if x, ok := _StatusValue[name]; ok {
		return x, nil
	}
	return Status(0), fmt.Errorf("%s is %w", name, ErrInvalidStatus)
}
