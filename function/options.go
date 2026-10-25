package function

// Option annotates an option of a function, usually used in variadic functions. An option applies a manipulation
// inplace to the given structure T and may fail when emitting an error.
type Option[T any] = func(*T) error

// ApplyOptions creates an options struct T and applies each Option to it. The final struct will be returned. If any
// Option returns an error, the first error fill be returned.
func ApplyOptions[T any](options []Option[T]) (T, error) {
	var opts T
	err := ApplyOptionsInto(options, &opts)
	if err != nil {
		return *new(T), err
	}
	return opts, nil
}

// ApplyOptionsInto works like ApplyOptions, but applies each Option into an already existing target options struct
// instead of allocating a new one.
func ApplyOptionsInto[T any](options []Option[T], opts *T) error {
	for _, option := range options {
		err := option(opts)
		if err != nil {
			return err
		}
	}
	return nil
}

// MustApplyOptionsInto works like ApplyOptionsInto but panics if an error occurs.
func MustApplyOptionsInto[T any](options []Option[T], opts *T) {
	err := ApplyOptionsInto(options, opts)
	if err != nil {
		panic(err)
	}
}

// MustApplyOptions works like ApplyOptions but panics if an error occurs.
func MustApplyOptions[T any](options []Option[T]) T {
	t, err := ApplyOptions[T](options)
	if err != nil {
		panic(err)
	}
	return t
}
