package stringutil

type StringSet map[string]struct{}

func NewStringSet(values []string) StringSet {
	result := StringSet{}
	for _, value := range values {
		result.Add(value)
	}
	return result
}

func (s *StringSet) Add(v string) {
	(*s)[v] = struct{}{}
}

func (s StringSet) Keys() []string {
	keys := make([]string, len(s))
	i := 0
	for k := range s {
		keys[i] = k
		i += 1
	}
	return keys
}

func (s StringSet) Contains(v string) bool {
	if _, ok := s[v]; ok {
		return true
	}
	return false
}
