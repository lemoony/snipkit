package utils

type StringSet map[string]struct{}

func (s StringSet) Keys() []string {
	keys := make([]string, len(s))
	i := 0
	for k := range s {
		keys[i] = k
		i += 1
	}
	return keys
}
