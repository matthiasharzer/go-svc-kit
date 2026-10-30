package iterutil_test

import (
	"errors"
	"iter"
	"testing"

	"github.com/matthiasharzer/service-kit/util/iterutil"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCollect(t *testing.T) {
	t.Run("collects all items from iterator", func(t *testing.T) {
		var iterator iter.Seq[string] = func(yield func(value string) bool) {
			yield("hellowo")
			yield("world")
		}
		collected := iterutil.Collect(iterator)

		require.Len(t, collected, 2)
		assert.Equal(t, "hellowo", collected[0])
		assert.Equal(t, "world", collected[1])
	})

	t.Run("returns empty slice for empty iterator", func(t *testing.T) {
		var iterator iter.Seq[string] = func(yield func(value string) bool) {}
		collected := iterutil.Collect(iterator)
		assert.Empty(t, collected)
	})
}

func TestCollect2(t *testing.T) {
	t.Run("collects all items from iterator", func(t *testing.T) {
		var iterator iter.Seq2[string, error] = func(yield func(value string, err error) bool) {
			yield("hellowo", nil)
			yield("world", nil)
		}
		collected, err := iterutil.Collect2(iterator)
		require.NoError(t, err)

		require.Len(t, collected, 2)
		assert.Equal(t, "hellowo", collected[0])
		assert.Equal(t, "world", collected[1])
	})

	t.Run("returns empty slice for empty iterator", func(t *testing.T) {
		var iterator iter.Seq2[string, error] = func(yield func(value string, err error) bool) {}
		collected, err := iterutil.Collect2(iterator)
		require.NoError(t, err)
		assert.Empty(t, collected)
	})

	t.Run("returns first error with nil slice", func(t *testing.T) {
		var iterator iter.Seq2[string, error] = func(yield func(value string, err error) bool) {
			yield("hellowo", nil)
			if !yield("", errors.New("iter error")) {
				return
			}
			yield("world", nil)
		}
		collected, err := iterutil.Collect2(iterator)
		require.Error(t, err)
		assert.EqualError(t, err, "iter error")
		assert.Empty(t, collected)
	})
}
