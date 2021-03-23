package jsonpatch

import (
	"strings"
)

const (
	separator = "/"
	wildcard  = "*"
	tilde     = "~"
)

// JSONPointer identifies a specific value within a JSON object specified in RFC 6901
type JSONPointer []string

// ParseJSONPointer converts a string into a JSONPointer
func ParseJSONPointer(str string) JSONPointer {
	return strings.Split(str, separator)
}

// String returns a string representation of a JSONPointer
func (p JSONPointer) String() string {
	return strings.Join(p, separator)
}

// Add adds an element to the JSONPointer
func (p JSONPointer) Add(elem string) JSONPointer {
	elem = strings.ReplaceAll(elem, tilde, "~0")
	elem = strings.ReplaceAll(elem, separator, "~1")
	return append(p, elem)
}

// Match matches a pattern which is a string JSONPointer which might also contains wildcards
func (p JSONPointer) Match(pattern string) bool {
	elements := strings.Split(pattern, separator)
	for i, element := range elements {
		if element == wildcard {
			continue
		} else if i >= len(p) || element != p[i] {
			return false
		}
	}

	return strings.HasSuffix(pattern, wildcard) || len(p) == len(elements)
}
