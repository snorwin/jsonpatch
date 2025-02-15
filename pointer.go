package jsonpatch

import (
	"slices"
	"strings"
)

// JSONPointer identifies a specific value within a JSON object specified in RFC 6901
type JSONPointer struct {
	path []string
	tags []string
}

func NewWithPrefix(prefix []string) JSONPointer {
	return JSONPointer{
		path: prefix,
		tags: []string{},
	}
}

// ParseJSONPointer converts a string into a JSONPointer
func ParseJSONPointer(str string) JSONPointer {
	return JSONPointer{
		path: strings.Split(str, separator),
		tags: []string{},
	}
}

// String returns a string representation of a JSONPointer
func (p JSONPointer) String() string {
	return strings.Join(p.path, separator)
}

// Add adds an element to the JSONPointer
func (p JSONPointer) Add(elem string) JSONPointer {
	elem = strings.ReplaceAll(elem, tilde, "~0")
	elem = strings.ReplaceAll(elem, separator, "~1")
	p.path = append(p.path, elem)
	return p
}

// Match matches a pattern which is a string JSONPointer which might also contains wildcards
func (p JSONPointer) Match(pattern string) bool {
	elements := strings.Split(pattern, separator)
	for i, element := range elements {
		if element == wildcard {
			continue
		} else if i >= len(p.path) || element != p.path[i] {
			return false
		}
	}

	return strings.HasSuffix(pattern, wildcard) || len(p.path) == len(elements)
}

// AddTags override tags referencing the JSONPointer
func (p JSONPointer) AddTags(tags []string) JSONPointer {
	p.tags = tags
	return p
}

// AddTags override tags referencing the JSONPointer
func (p JSONPointer) ShouldOmite() bool {
	return slices.Contains(p.tags, "omitempty")
}

// Prefix return a path as list to be used as prefix
func (p JSONPointer) Prefix() []string {
	return p.path
}
