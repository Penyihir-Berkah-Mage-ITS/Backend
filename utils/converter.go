package utils

import "strconv"

func StringToInteger(s string) int {
	converted, err := strconv.Atoi(s)
	if err != nil {
		return 0
	}
	return converted
}
