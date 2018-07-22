package helper

type TableFile interface {
	Sheets() []TableSheet
}

type TableSheet interface {
	GetValue(row, col int, isFloat bool) string

	Name() string

	IsFullRowEmpty(row int) bool
}
