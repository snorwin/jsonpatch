package jsonpatch

// Handler is the interfaces used by the walker to create patches
type Handler interface {
	// Add creates a JSONPatch with an 'add' operation and appends it to the patch list
	Add(pointer JSONPointer, modified interface{}) []JSONPatch

	// Remove creates a JSONPatch with an 'remove' operation and appends it to the patch list
	Remove(pointer JSONPointer, current interface{}) []JSONPatch

	// Replace creates a JSONPatch with an 'replace' operation and appends it to the patch list
	Replace(pointer JSONPointer, modified, current interface{}) []JSONPatch
}

// DefaultHandler implements the Handler
type DefaultHandler struct{}

// Add implements Handler
func (h *DefaultHandler) Add(pointer JSONPointer, value interface{}) []JSONPatch {
	// The 'add' operation either inserts a value into the array at the specified index or adds a new member to the object
	// NOTE: If the target location specifies an object member that does exist, that member's value is replaced
	return []JSONPatch{
		{
			Operation: "add",
			Path:      pointer.String(),
			Value:     value,
		},
	}
}

// Remove implements Handler
func (h *DefaultHandler) Remove(pointer JSONPointer, _ interface{}) []JSONPatch {
	// The 'remove' operation removes the value at the target location (specified by the pointer)
	return []JSONPatch{
		{
			Operation: "remove",
			Path:      pointer.String(),
		},
	}
}

// Replace implements Handler
func (h *DefaultHandler) Replace(pointer JSONPointer, value, _ interface{}) []JSONPatch {
	// The 'replace' operation replaces the value at the target location with a new value
	return []JSONPatch{
		{
			Operation: "replace",
			Path:      pointer.String(),
			Value:     value,
		},
	}
}
