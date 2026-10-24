package fstools

import (
	"fmt"
	"os"

	"github.com/matthiasharzer/go-svc-kit/function"
	"github.com/matthiasharzer/go-svc-kit/util/funcutils"
)

type temporaryFileOptions struct {
	FileEnding string
}

func WithFileEnding(fileEnding string) function.Option[temporaryFileOptions] {
	return func(t *temporaryFileOptions) error {
		t.FileEnding = fileEnding
		return nil
	}
}

func TemporaryFile(opts ...function.Option[temporaryFileOptions]) (string, func(), error) {
	options, err := function.ApplyOptions(opts)
	if err != nil {
		return "", nil, err
	}

	pattern := fmt.Sprintf("%s-*", getRandomPrefix())
	if options.FileEnding != "" {
		pattern += options.FileEnding
	}

	file, err := os.CreateTemp("", pattern)
	if err != nil {
		return "", nil, err
	}
	defer funcutils.LogError(file.Close, "failed to close temporary file")

	cleanup := func() {
		_ = os.Remove(file.Name())
	}

	return file.Name(), cleanup, nil
}
