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

// JSONPatchList is a list of JSONPatch
type JSONPatchList struct {
	list []JSONPatch
	raw  []byte
}

// Empty returns true if the JSONPatchList is empty
func (l JSONPatchList) Empty() bool {
	return l.Len() == 0
}

// Len returns the length of the JSONPatchList
func (l JSONPatchList) Len() int {
	return len(l.list)
}

// String returns the encoded JSON string of the JSONPatchList
func (l JSONPatchList) String() string {
	return string(l.raw)
}

// Raw returns the raw encoded JSON of the JSONPatchList
func (l JSONPatchList) Raw() []byte {
	return l.raw
}

// List returns a copy of the underlying JSONPatch slice
func (l JSONPatchList) List() []JSONPatch {
	ret := make([]JSONPatch, l.Len())

	for i, patch := range l.list {
		ret[i] = patch
	}

	return ret
}

// CreateJSONPatch compares two JSON data structures and creates a JSONPatch according to RFC 6902
func CreateJSONPatch(modified, current interface{}, options ...Option) (JSONPatchList, error) {
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
		return JSONPatchList{}, err
	}

	list := w.patchList
	if len(list) == 0 {
		return JSONPatchList{}, nil
	}
	raw, err := json.Marshal(list)

	return JSONPatchList{list: list, raw: raw}, err
}

// CreateThreeWayJSONPatch compares three JSON data structures and creates a three-way JSONPatch according to RFC 6902
func CreateThreeWayJSONPatch(modified, current, original interface{}, options ...Option) (JSONPatchList, error) {
	var list []JSONPatch

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

	// compare modified with current and only keep addition and changes
	if err := w.walk(reflect.ValueOf(modified), reflect.ValueOf(current), w.prefix); err != nil {
		return JSONPatchList{}, err
	}
	for _, patch := range w.patchList {
		if patch.Operation != "remove" {
			list = append(list, patch)
		}
	}

	// reset walker
	w.patchList = []JSONPatch{}

	// compare modified with original and only keep deletions
	if err := w.walk(reflect.ValueOf(modified), reflect.ValueOf(original), w.prefix); err != nil {
		return JSONPatchList{}, err
	}
	for _, patch := range w.patchList {
		if patch.Operation == "remove" {
			list = append(list, patch)
		}
	}

	if len(list) == 0 {
		return JSONPatchList{}, nil
	}
	raw, err := json.Marshal(list)

	return JSONPatchList{list: list, raw: raw}, err
}
