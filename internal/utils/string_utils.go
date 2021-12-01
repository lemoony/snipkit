package utils

func StringOrDefault(val string, defaultStr string) string {
	if val == "" {
		return defaultStr
	}
	return val
}
