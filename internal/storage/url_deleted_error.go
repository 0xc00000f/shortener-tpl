package storage

type URLDeletedError struct{}

func (fde URLDeletedError) Error() string {
	return "this url has been already deleted"
}
