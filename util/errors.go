package util

type DostupErrorKind int

const (
	NotFound DostupErrorKind = iota
	Conflict
	Unavailable
	InvalidParameter
)

type DostupError struct {
	kind    DostupErrorKind
	err     string
	message string
}

func DErr(kind DostupErrorKind, err string) *DostupError {
	return &DostupError{
		kind: kind,
		err:  err,
	}
}
func (de *DostupError) Error() string {
	return de.err
}
func (de *DostupError) SetMessage(message string) *DostupError {
	de.message = message
	return de
}
func (de *DostupError) Message() string {
	return de.message
}
func (de *DostupError) Kind() DostupErrorKind {
	return de.kind
}
