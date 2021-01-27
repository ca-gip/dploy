package utils

import (
	"sort"
	"strings"
)

type emptyType = struct{}

var empty = emptyType{}

type Set struct {
	m map[string]struct{}
}

func NewSet() *Set {
	s := &Set{}
	return s.initIfEmpty()
}

func NewSetFromSlice(elements ...string) *Set {
	s := &Set{}
	return s.initIfEmpty().Concat(elements)
}

func (s *Set) initIfEmpty() *Set {
	if s.m == nil {
		s.m = make(map[string]struct{})
	}
	return s
}

func (s *Set) Add(value string) *Set {
	s.initIfEmpty().m[value] = empty
	return s
}

func (s *Set) Concat(values []string) *Set {
	s.initIfEmpty()
	for _, value := range values {
		s.Add(value)
	}
	return s
}

func (s *Set) Remove(value string) *Set {
	s.initIfEmpty()
	delete(s.m, value)
	return s
}

// Return a sorted list of values
func (s *Set) List() (list []string) {
	for key := range s.initIfEmpty().m {
		list = append(list, key)
	}
	sort.Strings(list)
	return
}

func (s *Set) Contains(value string) bool {
	_, c := s.initIfEmpty().m[value]
	return c
}

// Read Struct Field to return the associated Tag name
// It's used to read string or slice of string due to different notation in yaml
func (s *Set) UnmarshalYAML(unmarshal func(i interface{}) error) (err error) {
	var sliceReceiver []string
	var stringReceiver string

	if err = unmarshal(&sliceReceiver); err == nil {
		s.Concat(sliceReceiver)
		return nil
	} else if err = unmarshal(&stringReceiver); err == nil {
		s.Concat(strings.Split(stringReceiver, ","))
		return nil
	}
	return err
}
