package utils

var exists = struct{}{}

type Set struct {
	m map[string]struct{}
}

func NewSet() *Set {
	s := &Set{}
	s.m = make(map[string]struct{})
	return s
}

func (s *Set) Add(value string) *Set {
	s.m[value] = exists
	return s
}

func (s *Set) Concat(values []string) *Set {
	for _, value := range values {
		s.Add(value)
	}
	return s
}

func (s *Set) Remove(value string) *Set {
	delete(s.m, value)
	return s
}

func (s *Set) List() (list []string) {
	for key, _ := range s.m {
		list = append(list, key)
	}
	return
}

func (s *Set) Contains(value string) bool {
	_, c := s.m[value]
	return c
}
