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
	EndWith   = "$="
	Contains  = "~="
	StartWith = "^="
)

var (
	AllowedOperators   = []string{Equal, NotEqual, StartWith, Contains, EndWith}
	FiltersRe          = regexp.MustCompile("(\\w*)(==|!=|$=|~=|^=)(\\w*)")
	FilterRegex        = regexp.MustCompile("(\\w*)(\\W*)(\\w*)")
	FilterCompletionRe = regexp.MustCompile("(\\w*)(==|!=|$=|~=|^=)(\\w*)([,]|)")
)

type Filter struct {
	Key   string
	Op    string
	Value string
}

type Filters []Filter

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
	case Equal:
		return strings.EqualFold(f.Value, against)
	case NotEqual:
		return !strings.EqualFold(f.Value, against)
	case EndWith:
		return strings.HasSuffix(f.Value, against)
	case Contains:
		return strings.Contains(f.Value, against)
	case StartWith:
		return strings.HasPrefix(f.Value, against)
	default:
		log.Fatalf("Unsuported filter operation %s", f.Op)
		return false
	}
}
