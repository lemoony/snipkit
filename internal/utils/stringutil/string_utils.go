package stringutil

import "strings"

func StringOrDefault(val string, defaultStr string) string {
	if val == "" {
		return defaultStr
	}
	return val
}

// SplitWithEscape behaves like strings.Split but supports defining an escape character. E.g, when using ',' as split
// separator and '\' es escape code, the string "1\,2,3," will not be split into ["1, 2", "3"].
func SplitWithEscape(s string, split uint8, escape uint8, trim bool) []string {
	var splitIndices []int
	for i := range s {
		if s[i] == split && (i == 0 || s[i-1] != escape) {
			splitIndices = append(splitIndices, i)
		}
	}

	var result []string

	for i, lowerSplitIndex := 0, 0; i <= len(splitIndices); i++ {
		if i > 0 {
			lowerSplitIndex = splitIndices[i-1] + 1
		}

		upperIndex := len(s)
		if i < len(splitIndices) {
			upperIndex = splitIndices[i]
		}

		sp := s[lowerSplitIndex:upperIndex]
		sp = strings.ReplaceAll(sp, "\\,", ",")
		if trim {
			sp = strings.TrimSpace(sp)
		}

		if sp != "" {
			result = append(result, sp)
		}
	}

	return result
}

func FirstNotEmpty(values ...string) string {
	for _, value := range values {
		if value != "" {
			return value
		}
	}
	return ""
}
