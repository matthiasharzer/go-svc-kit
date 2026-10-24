package function

type Option[T any] = func(*T) error

func ApplyOptionsInto[T any](options []Option[T], opts *T) error {
	for _, option := range options {
		err := option(opts)
		if err != nil {
			return err
		}
	}
	return nil
}

func ApplyOptions[T any](options []Option[T]) (T, error) {
	var opts T
	err := ApplyOptionsInto(options, &opts)
	if err != nil {
		return *new(T), err
	}
	return opts, nil
}

func MustApplyOptionsInto[T any](options []Option[T], opts *T) {
	err := ApplyOptionsInto(options, opts)
	if err != nil {
		panic(err)
	}
}

func MustApplyOptions[T any](options []Option[T]) T {
	t, err := ApplyOptions[T](options)
	if err != nil {
		panic(err)
	}
	return t
}
