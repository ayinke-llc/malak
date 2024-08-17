package malak

type malakError string

func (m malakError) Error() string { return string(m) }
