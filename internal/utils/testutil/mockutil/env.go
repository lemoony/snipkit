package mockutil

import "os"

func MockAPIKeyEnv(key, value string) func() {
	originalValue := os.Getenv(key)
	_ = os.Setenv(key, value)
	return func() {
		_ = os.Setenv(key, originalValue)
	}
}
