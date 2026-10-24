package function_test

import (
	"errors"
	"testing"
	"time"

	"github.com/matthiasharzer/go-svc-kit/function"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type testOptions struct {
	someValue    string
	someDuration time.Duration
}

func withSomeValue(value string) function.Option[testOptions] {
	return func(t *testOptions) error {
		t.someValue = value
		return nil
	}
}

func withPositiveDuration(value time.Duration) function.Option[testOptions] {
	return func(t *testOptions) error {
		if value < 0 {
			return errors.New("negative duration")
		}
		t.someDuration = value
		return nil
	}
}

func TestApplyOptions(t *testing.T) {
	t.Run("applies all options", func(t *testing.T) {
		opts := []function.Option[testOptions]{
			withSomeValue("hello"),
			withPositiveDuration(1 * time.Second),
		}

		appliedOptions, err := function.ApplyOptions(opts)
		require.NoError(t, err)

		assert.Equal(t, appliedOptions.someValue, "hello")
		assert.Equal(t, appliedOptions.someDuration, 1*time.Second)
	})

	t.Run("returns empty options", func(t *testing.T) {
		appliedOptions, err := function.ApplyOptions[testOptions](nil)
		require.NoError(t, err)
		assert.Equal(t, appliedOptions.someValue, "")
		assert.Equal(t, appliedOptions.someDuration, time.Duration(0))
	})

	t.Run("returns error if option is invalid", func(t *testing.T) {
		opts := []function.Option[testOptions]{
			withSomeValue("hello"),
			withPositiveDuration(-1 * time.Second),
		}

		_, err := function.ApplyOptions(opts)
		require.Error(t, err)
		assert.EqualError(t, err, "negative duration")
	})
}

func TestApplyOptionsTo(t *testing.T) {
	t.Run("applies all options into existing target", func(t *testing.T) {
		opts := []function.Option[testOptions]{
			withSomeValue("hello"),
			withPositiveDuration(1 * time.Second),
		}

		var target testOptions

		err := function.ApplyOptionsInto(opts, &target)
		require.NoError(t, err)

		assert.Equal(t, target.someValue, "hello")
		assert.Equal(t, target.someDuration, 1*time.Second)
	})

	t.Run("substitutes nothing if options are empty", func(t *testing.T) {
		var target testOptions
		err := function.ApplyOptionsInto(nil, &target)
		require.NoError(t, err)
		assert.Equal(t, target.someValue, "")
		assert.Equal(t, target.someDuration, time.Duration(0))
	})

	t.Run("returns error if option is invalid", func(t *testing.T) {
		opts := []function.Option[testOptions]{
			withSomeValue("hello"),
			withPositiveDuration(-1 * time.Second),
		}

		var target testOptions
		err := function.ApplyOptionsInto(opts, &target)
		require.Error(t, err)
		assert.EqualError(t, err, "negative duration")
	})
}

func TestMustApplyOptions(t *testing.T) {
	t.Run("panics if option is invalid", func(t *testing.T) {
		opts := []function.Option[testOptions]{
			withSomeValue("hello"),
			withPositiveDuration(-1 * time.Second),
		}
		defer func() {
			if r := recover(); r == nil {
				assert.Fail(t, "function did not panic")
			}
		}()
		_ = function.MustApplyOptions(opts)
	})
}

func TestMustApplyOptionsInto(t *testing.T) {
	t.Run("panics if option is invalid", func(t *testing.T) {
		opts := []function.Option[testOptions]{
			withSomeValue("hello"),
			withPositiveDuration(-1 * time.Second),
		}
		defer func() {
			if r := recover(); r == nil {
				assert.Fail(t, "function did not panic")
			}
		}()
		var target testOptions
		function.MustApplyOptionsInto(opts, &target)
	})
}
