package malak

type MalakError string

func (m MalakError) Error() string { return string(m) }
