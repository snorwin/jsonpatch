package jsonpatch

import (
	"encoding/json"
	"reflect"
)

// JSONPatch format is specified in RFC 6902
type JSONPatch struct {
	Operation string      `json:"op"`
	Path      string      `json:"path"`
	Value     interface{} `json:"value,omitempty"`
}

// Patch is a byte encoded JSON patch
type Patch []byte

// Empty returns true if the Patch is empty
func (p Patch) Empty() bool {
	return len(p) == 0
}

// String converts the Patch to a string
func (p Patch) String() string {
	return string(p)
}

// CreateJSONPatch compares two JSON data structures and creates a JSONPatch according to RFC 6902
func CreateJSONPatch(modified, current interface{}, options ...Option) (Patch, int, error) {
	// create a new walker
	w := &walker{
		handler:   &DefaultHandler{},
		predicate: Funcs{},
		prefix:    []string{""},
	}

	// apply options to the walker
	for _, apply := range options {
		apply(w)
	}

	if err := w.walk(reflect.ValueOf(modified), reflect.ValueOf(current), w.prefix); err != nil {
		return []byte{}, 0, err
	}

	patches := w.patchList
	if len(patches) == 0 {
		return []byte{}, 0, nil
	}
	p, err := json.Marshal(patches)

	return p, len(patches), err
}
