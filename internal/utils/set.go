package utils

import (
	"strings"
)

type emptyType = struct{}

var empty = emptyType{}

type Set struct {
	m map[string]struct{}
}

func NewSet() *Set {
	s := &Set{}
	s.m = make(map[string]struct{})
	return s
}

func NewSetFromSlice(elements ...string) *Set {
	s := &Set{}
	s.m = make(map[string]struct{})
	return s.Concat(elements)
}

func (s *Set) Add(value string) *Set {
	s.m[value] = empty
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

// Read Struct Field to return the associated Tag name
// It's used to read string or slice of string due to different notation in yaml
func (s *Set) UnmarshalYAML(unmarshal func(i interface{}) error) (err error) {
	var sliceReceiver []string
	var stringReceiver string

	if err = unmarshal(&sliceReceiver); err == nil {
		s.m = make(map[string]emptyType, len(sliceReceiver))
		for _, item := range sliceReceiver {
			s.m[item] = emptyType{}
		}
		return nil
	} else if err = unmarshal(&stringReceiver); err == nil {
		strSplits := strings.Split(stringReceiver, ",")
		s.m = make(map[string]emptyType, len(sliceReceiver))
		for _, item := range strSplits {
			s.m[item] = emptyType{}
		}
		return nil
	}
	return err
}
