package jsonpatch

import (
	"fmt"
	"reflect"
	"sort"
	"strconv"
	"strings"
)

const (
	jsonTag = "json"
)

type walker struct {
	predicate     Predicate
	handler       Handler
	prefix        []string
	patchList     []JSONPatch
	ignoredSlices []IgnorePattern
}

// walk recursively processes the modified and current JSON data structures simultaneously and in every step it compares
// the value of them with each other
func (w *walker) walk(modified, current reflect.Value, pointer JSONPointer) error {
	// the data structures of both JSON objects must be identical
	if modified.Kind() != current.Kind() {
		return fmt.Errorf("kind does not match at: %s modified: %s current: %s", pointer, modified.Kind(), current.Kind())
	}
	switch modified.Kind() {
	case reflect.Struct:
		return w.processStruct(modified, current, pointer)
	case reflect.Ptr:
		return w.processPtr(modified, current, pointer)
	case reflect.Slice:
		return w.processSlice(modified, current, pointer)
	case reflect.Map:
		return w.processMap(modified, current, pointer)
	case reflect.Interface:
		return w.processInterface(modified, current, pointer)
	case reflect.String:
		if modified.String() != current.String() {
			if modified.String() == "" {
				w.remove(pointer, current.String())
			} else if current.String() == "" {
				w.add(pointer, modified.String())
			} else {
				w.replace(pointer, modified.String(), current.String())
			}
		}
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		if modified.Int() != current.Int() {
			w.replace(pointer, modified.Int(), current.Int())
		}
	case reflect.Uint, reflect.Uintptr, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		if modified.Uint() != current.Uint() {
			w.replace(pointer, modified.Uint(), current.Uint())
		}
	case reflect.Bool:
		if modified.Bool() != current.Bool() {
			w.replace(pointer, modified.Bool(), current.Bool())
		}
	case reflect.Float32, reflect.Float64:
		if modified.Float() != current.Float() {
			w.replace(pointer, modified.Float(), current.Float())
		}
	case reflect.Invalid:
		// undefined interfaces are ignored for now
		return nil
	default:
		return fmt.Errorf("unsupported kind: %s at: %s", modified.Kind(), pointer)
	}

	return nil
}

// processInterface processes reflect.Interface values
func (w *walker) processInterface(modified reflect.Value, current reflect.Value, pointer JSONPointer) error {
	// extract the value form the interface and try to process it further
	if err := w.walk(reflect.ValueOf(modified.Interface()), reflect.ValueOf(current.Interface()), pointer); err != nil {
		return err
	}

	return nil
}

// processMap processes reflect.Map values
func (w *walker) processMap(modified reflect.Value, current reflect.Value, pointer JSONPointer) error {
	// NOTE: currently only map[string]interface{} are supported
	if len(modified.MapKeys()) > 0 && len(current.MapKeys()) == 0 {
		w.add(pointer, modified.Interface())
	} else {
		it := modified.MapRange()
		for it.Next() {
			key := it.Key()
			if key.Kind() != reflect.String {
				return fmt.Errorf("only strings are supported as map keys but was: %s at: %s", key.Kind(), pointer)
			}

			val1 := it.Value()
			val2 := current.MapIndex(key)
			if val2.Kind() == reflect.Invalid {
				w.add(pointer.Add(key.String()), val1.Interface())
			} else {
				if err := w.walk(val1, val2, pointer.Add(key.String())); err != nil {
					return err
				}
			}
		}
		it = current.MapRange()
		for it.Next() {
			key := it.Key()
			if key.Kind() != reflect.String {
				return fmt.Errorf("only strings are supported as map keys but was: %s at: %s", key.Kind(), pointer)
			}

			val1 := modified.MapIndex(key)
			val2 := it.Value()
			if val1.Kind() == reflect.Invalid {
				w.remove(pointer.Add(key.String()), val2.Interface())
			}
		}
	}

	return nil
}

// processSlice processes reflect.Slice values
func (w *walker) processSlice(modified reflect.Value, current reflect.Value, pointer JSONPointer) error {
	if !w.predicate.Replace(pointer, modified.Interface(), current.Interface()) {
		return nil
	}

	if modified.Len() > 0 && current.Len() == 0 {
		w.add(pointer, modified.Interface())
	} else {
		var ignoreSliceOrder bool
		var patchSliceJSONField string

		// look for a slice tag which pattern matches the pointer
		for _, ignore := range w.ignoredSlices {
			if pointer.Match(ignore.Pattern) {
				ignoreSliceOrder = true
				patchSliceJSONField = ignore.JSONField
				break
			}
		}

		if ignoreSliceOrder {
			fieldName := jsonFieldNameToFieldName(modified.Type().Elem(), patchSliceJSONField)

			// maps the modified slice elements with the patchSliceKey to their index
			idxMap1 := map[string]int{}
			for j := 0; j < modified.Len(); j++ {
				fieldValue := extractIgnoreSliceOrderMatchValue(modified.Index(j), fieldName)
				if _, ok := idxMap1[fieldValue]; ok {
					return fmt.Errorf("ignore slice order failed at %s due to unique match field constraint, duplicated value: %s", pointer, fieldValue)
				}
				idxMap1[fieldValue] = j
			}

			// maps the current slice elements with the patchSliceKey to their index
			idxMap2 := map[string]int{}
			for j := 0; j < current.Len(); j++ {
				fieldValue := extractIgnoreSliceOrderMatchValue(current.Index(j), fieldName)
				if _, ok := idxMap2[fieldValue]; ok {
					return fmt.Errorf("ignore slice order failed at %s due to unique match field constraint, duplicated value: %s", pointer, fieldValue)
				}
				idxMap2[fieldValue] = j
			}

			// IMPORTANT: the order of the patches matters, because and add or delete will change the index of your
			// elements. Therefore, elements are only added at the end of the slice and all elements are updated
			// before any element is deleted.

			// iterate through the list of modified slice elements in order to identify updated or added elements
			idxMax := current.Len()
			for k, idx1 := range idxMap1 {
				idx2, ok := idxMap2[k]
				if ok {
					if err := w.walk(modified.Index(idx1), current.Index(idx2), pointer.Add(strconv.Itoa(idx2))); err != nil {
						return err
					}
				} else {
					if ok := w.add(pointer.Add(strconv.Itoa(idxMax)), modified.Index(idx1).Interface()); ok {
						idxMax++
					}
				}
			}

			// IMPORTANT: deleting must be done in reverse order

			// iterate through the list of current slice elements in order to identify deleted elements
			var deleted []int
			for k, idx2 := range idxMap2 {
				if _, ok := idxMap1[k]; !ok {
					deleted = append(deleted, idx2)
				}
			}
			sort.Ints(deleted)
			for j := len(deleted) - 1; j >= 0; j-- {
				idx2 := deleted[j]
				w.remove(pointer.Add(strconv.Itoa(idx2)), current.Index(idx2).Interface())
			}
		} else {
			// iterate through both slices and update their elements until on of them is completely processed
			for j := 0; j < modified.Len() && j < current.Len(); j++ {
				if err := w.walk(modified.Index(j), current.Index(j), pointer.Add(strconv.Itoa(j))); err != nil {
					return err
				}
			}
			if modified.Len() > current.Len() {
				// add the remaining elements of the modified slice
				idx := current.Len()
				for j := current.Len(); j < modified.Len(); j++ {
					if ok := w.add(pointer.Add(strconv.Itoa(idx)), modified.Index(j).Interface()); ok {
						idx++
					}
				}
			} else if modified.Len() < current.Len() {
				// delete the remaining elements of the current slice
				// IMPORTANT: deleting must be done in reverse order
				for j := current.Len() - 1; j >= modified.Len(); j-- {
					w.remove(pointer.Add(strconv.Itoa(j)), current.Index(j).Interface())
				}
			}
		}
	}

	return nil
}

// processPtr processes reflect.Ptr values
func (w *walker) processPtr(modified reflect.Value, current reflect.Value, pointer JSONPointer) error {
	if !modified.IsNil() && !current.IsNil() {
		// the values of the pointers will be processed in a next step
		if err := w.walk(modified.Elem(), current.Elem(), pointer); err != nil {
			return err
		}
	} else if !modified.IsNil() {
		w.add(pointer, modified.Elem().Interface())
	} else if !current.IsNil() {
		w.remove(pointer, current.Elem().Interface())
	}

	return nil
}

// processStruct processes reflect.Struct values
func (w *walker) processStruct(modified, current reflect.Value, pointer JSONPointer) error {
	if !w.predicate.Replace(pointer, modified.Interface(), current.Interface()) {
		return nil
	}

	// process all struct fields, the order of the fields of the  modified and current JSON object is identical because their types match
	for j := 0; j < modified.NumField(); j++ {
		tag := strings.Split(modified.Type().Field(j).Tag.Get(jsonTag), ",")[0]
		if tag == "" || tag == "_" || !modified.Field(j).CanInterface() {
			// struct fields without a JSON tag set or unexported fields are ignored
			continue
		}
		// process the child's value of the modified and current JSON in a next step
		if err := w.walk(modified.Field(j), current.Field(j), pointer.Add(tag)); err != nil {
			return err
		}
	}

	return nil
}

// extractIgnoreSliceOrderMatchValue extracts the value which is used to match the modified and current values to ignore the slice order
func extractIgnoreSliceOrderMatchValue(value reflect.Value, fieldName string) string {
	switch value.Kind() {
	case reflect.Struct:
		return extractIgnoreSliceOrderMatchValue(value.FieldByName(fieldName), "")
	case reflect.Ptr:
		if !value.IsNil() {
			return extractIgnoreSliceOrderMatchValue(value.Elem(), fieldName)
		}
		return ""
	case reflect.String:
		return value.String()
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return strconv.FormatInt(value.Int(), 10)
	case reflect.Uint, reflect.Uintptr, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return strconv.FormatUint(value.Uint(), 10)
	case reflect.Bool:
		return strconv.FormatBool(value.Bool())
	}
	return ""
}

// jsonFieldNameToFieldName retrieves the actual Go field name for the JSON field name
func jsonFieldNameToFieldName(t reflect.Type, jsonFieldName string) string {
	switch t.Kind() {
	case reflect.Struct:
		for i := 0; i < t.NumField(); i++ {
			if strings.Split(t.Field(i).Tag.Get(jsonTag), ",")[0] == jsonFieldName {
				return t.Field(i).Name
			}
		}
	case reflect.Ptr:
		return jsonFieldNameToFieldName(t.Elem(), jsonFieldName)
	}

	return ""
}

// add adds an add JSON patch by checking the Predicate first and using the Handler to generate it
func (w *walker) add(pointer JSONPointer, modified interface{}) bool {
	if w.predicate != nil && !w.predicate.Add(pointer, modified) {
		return false
	}
	w.patchList = append(w.patchList, w.handler.Add(pointer, modified)...)

	return true
}

// replace adds a replace JSON patch by checking the Predicate first and using the Handler to generate it
func (w *walker) replace(pointer JSONPointer, modified, current interface{}) bool {
	if w.predicate != nil && !w.predicate.Replace(pointer, modified, current) {
		return false
	}
	w.patchList = append(w.patchList, w.handler.Replace(pointer, modified, current)...)

	return true
}

// remove adds a remove JSON patch by checking the Predicate first and using the Handler to generate it
func (w *walker) remove(pointer JSONPointer, current interface{}) bool {
	if w.predicate != nil && !w.predicate.Remove(pointer, current) {
		return false
	}
	w.patchList = append(w.patchList, w.handler.Remove(pointer, current)...)

	return true
}
