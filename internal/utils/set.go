package utils

import (
	"fmt"
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
// omitempty can be specified
func (s *Set) UnmarshalYAML(unmarshal func(i interface{}) error) (err error) {
	var tmpSlice []string
	var tmpString string

	if err = unmarshal(&tmpSlice); err == nil {
		s.m = make(map[string]emptyType, len(tmpSlice))

		for _, v := range tmpSlice {
			s.m[v] = emptyType{} // add 1 to k as it is starting at base 0
		}
		fmt.Println("slice string read", s)
		return nil
	} else if err = unmarshal(&tmpString); err == nil {
		strSplits := strings.Split(tmpString, ",")
		s.m = make(map[string]emptyType, len(tmpSlice))

		for _, v := range strSplits {
			s.m[v] = emptyType{} // add 1 to k as it is starting at base 0
		}
		fmt.Println("string", s)
		return nil
	}
	return err
}
