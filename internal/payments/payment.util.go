package payments

import "strconv"

func parsePositiveInt(value string, defaultValue int) int {
	if value == "" {
		return defaultValue
	}

	parsedValue, err := strconv.Atoi(value)
	if err != nil || parsedValue < 1 {
		return defaultValue
	}

	return parsedValue
}
