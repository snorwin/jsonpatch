package jsonpatch

// Predicate filters patches
type Predicate interface {
	// Add returns true if the object should not be added in the patch
	Add(pointer JSONPointer, modified interface{}) bool

	// Remove returns true if the object should not be deleted in the patch
	Remove(pointer JSONPointer, current interface{}) bool

	// Replace returns true if the objects should not be updated in the patch - this will stop the recursive processing of those objects
	Replace(pointer JSONPointer, modified, current interface{}) bool
}

// Funcs is a function that implements Predicate
type Funcs struct {
	// Add returns true if the object should not be added in the patch
	AddFunc func(pointer JSONPointer, modified interface{}) bool

	// Remove returns true if the object should not be deleted in the patch
	RemoveFunc func(pointer JSONPointer, current interface{}) bool

	// Replace returns true if the objects should not be updated in the patch - this will stop the recursive processing of those objects
	ReplaceFunc func(pointer JSONPointer, modified, current interface{}) bool
}

// Add implements Predicate
func (p Funcs) Add(pointer JSONPointer, modified interface{}) bool {
	if p.AddFunc != nil {
		return p.AddFunc(pointer, modified)
	}

	return true
}

// Remove implements Predicate
func (p Funcs) Remove(pointer JSONPointer, current interface{}) bool {
	if p.RemoveFunc != nil {
		return p.RemoveFunc(pointer, current)
	}

	return true
}

// Replace implements Predicate
func (p Funcs) Replace(pointer JSONPointer, modified, current interface{}) bool {
	if p.ReplaceFunc != nil {
		return p.ReplaceFunc(pointer, modified, current)
	}

	return true
}
