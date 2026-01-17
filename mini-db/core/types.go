package core

import "fmt"

type DataType string

const (
	IntType  DataType = "INT"
	TextType DataType = "TEXT"
)

func ParseDataType(s string) (DataType, error) {
	switch s {
	case "INT":
		return IntType, nil
	case "TEXT":
		return TextType, nil
	default:
		return "", fmt.Errorf("unknown data type: %s", s)
	}
}
