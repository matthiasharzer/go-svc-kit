package itertools

import "iter"

// Collect retrieves all elements from the given iter.Seq iterator and returning them as a slice.
func Collect[T any](iterator iter.Seq[T]) []T {
	var result []T
	for value := range iterator {
		result = append(result, value)
	}
	return result
}

// Collect2 returns all items from an iter.Seq2 which second element if of type error. The collected values (first type)
// will be returned as a slice. If any error is returned from the iterator, it will be returned and iteration canceled.
func Collect2[T any](iterator iter.Seq2[T, error]) ([]T, error) {
	var result []T
	for value, err := range iterator {
		if err != nil {
			return nil, err
		}
		result = append(result, value)
	}
	return result, nil
}
