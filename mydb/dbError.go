package mydb

type dbError struct {
	text string
}

func newDbError(text string) *dbError {
	return &dbError{text: "DataBaseError: " + text}
}

func (e dbError) Error() string {
	return e.text
}
