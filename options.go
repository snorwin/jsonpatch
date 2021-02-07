package jsonpatch

// Option allow to configure the walker instance
type Option func(r *walker)

// WithPredicate set a patch Predicate for the walker. This can be used to filter or validate the patch creation
func WithPredicate(predicate Predicate) Option {
	return func(w *walker) {
		w.predicate = predicate
	}
}

// WithHandler set a patch Handler for the walker. This can be used to customize the patch creation
func WithHandler(handler Handler) Option {
	return func(w *walker) {
		w.handler = handler
	}
}

// WithPrefix is used to specify a prefix if only a sub part of JSON structure needs to be patched
func WithPrefix(prefix []string) Option {
	return func(w *walker) {
		if len(prefix) > 0 && prefix[0] == "" {
			w.prefix = append(w.prefix, prefix[1:]...)
		} else {
			w.prefix = append(w.prefix, prefix...)
		}
	}
}

// IgnoreSliceOrder will ignore the order of all slices of built-in types during the walk and will use instead the value
// itself in order to compare  the current and modified JSON.
// NOTE: ignoring order only works if the elements in each slice are unique
func IgnoreSliceOrder() Option {
	return func(w *walker) {
		w.ignoredSlices = append(w.ignoredSlices, IgnorePattern{Pattern: "*"})
	}
}

// IgnoreSliceOrderWithPattern will ignore the order of slices which paths match the pattern during the walk
// and will use instead the value in order to compare  the current and modified JSON.
// (For structs (and pointers of structs) the value of the JSON field specified in IgnorePattern is used for comparison.)
// NOTE: ignoring order only works if the elements in each slice are unique
func IgnoreSliceOrderWithPattern(slices []IgnorePattern) Option {
	return func(w *walker) {
		w.ignoredSlices = append(slices, w.ignoredSlices...)
	}
}

// IgnorePattern specifies a JSONPointer Pattern and a an optional JSONField
type IgnorePattern struct {
	Pattern   string
	JSONField string
}
