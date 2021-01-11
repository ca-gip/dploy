package services

import (
	"fmt"
	"log"
	"regexp"
	"strings"
)

const (
	Equal     = "=="
	NotEqual  = "!="
	StartWith = "$="
	Contains  = "~="
	EndWith   = "^="
)

var (
	AllowedOperators = []string{Equal, NotEqual, StartWith, Contains, EndWith}
	FiltersRe        = regexp.MustCompile("(\\w*)(==|!=|$=|~=|^=)(\\w*)")
	FilterRegex      = regexp.MustCompile("(\\w*)(\\W*)(\\w*)")
)

type Filter struct {
	Key   string
	Op    string
	Value string
}

func ParseFilter(filter string) (key, op, value string) {
	result := FilterRegex.FindStringSubmatch(filter)
	return result[1], result[2], result[3]
}

func ParseFilterArgsFromSlice(rawFilters []string) (filters []Filter) {
	for _, rawFilter := range rawFilters {
		result := FilterRegex.FindStringSubmatch(rawFilter)
		filters = append(filters, Filter{
			Key:   result[1],
			Op:    result[2],
			Value: result[3],
		})
	}
	return
}

func ParseFilterArgsFromString(raw string) (filters []Filter) {
	for _, filter := range FiltersRe.FindAllStringSubmatch(raw, -1) {
		fmt.Println(filter)
		filters = append(filters, Filter{
			Key:   filter[1],
			Op:    filter[2],
			Value: filter[3],
		})
	}
	return
}

func (f Filter) GetRaw() string {
	return fmt.Sprintf("%s%s%s", f.Key, f.Op, f.Value)
}

func (f Filter) Eval(against string) bool {
	switch f.Op {
	case "==":
		return strings.EqualFold(f.Value, against)
	case "!=":
		return !strings.EqualFold(f.Value, against)
	case "$=":
		return strings.HasSuffix(f.Value, against)
	case "~=":
		return strings.Contains(f.Value, against)
	case "^=":
		return strings.HasPrefix(f.Value, against)
	default:
		log.Fatalf("Unsuported filter operation %s", f.Op)
		return false
	}
}
