package utils

var exists = struct{}{}

type set struct {
	m map[string]struct{}
}

func NewSet() *set {
	s := &set{}
	s.m = make(map[string]struct{})
	return s
}

func (s *set) Add(value string) *set {
	s.m[value] = exists
	return s
}

func (s *set) Concat(values []string) *set {
	for _, value := range values {
		s.Add(value)
	}
	return s
}

func (s *set) Remove(value string) *set {
	delete(s.m, value)
	return s
}

func (s *set) List() (list []string) {
	for key, _ := range s.m {
		list = append(list, key)
	}
	return
}

func (s *set) Contains(value string) bool {
	_, c := s.m[value]
	return c
}
