package malak

import "strings"

type MalakError string

func (m MalakError) Error() string { return string(m) }

func IsDuplicateUniqueError(e error) bool {
	return strings.Contains(e.Error(), "duplicate key value violates unique constraint")
}
