package utils

import (
	"fmt"
	"strconv"
)

func ConvertStr(value string) (int, error) {
	if value == "" {
		return 0, fmt.Errorf("value is empty")

	}
	converted_value, err := strconv.Atoi(value)
	if err != nil {
		return 0, fmt.Errorf("invalid value: %w", err)
	}

	return converted_value, nil
}
