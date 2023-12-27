package utils

func StringOrDefault(str string, def string) string {
	if len(str) > 0 {
		return str
	}

	return def
}
