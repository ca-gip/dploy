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

func (s *set) Add(value string) {
	s.m[value] = exists
}

func (s *set) Remove(value string) {
	delete(s.m, value)
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
