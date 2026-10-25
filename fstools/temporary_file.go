package fstools

import (
	"fmt"
	"os"

	"github.com/matthiasharzer/service-kit/function"
)

type temporaryFileOptions struct {
	FileEnding string
}

// WithFileEnding configures the file ending to use when creating a temporary file. This is a pure suffix, thus it
// does not include a dot. E.g. when you want the file to end in ".txt", you need to use ".txt" as the file ending.
func WithFileEnding(fileEnding string) function.Option[temporaryFileOptions] {
	return func(t *temporaryFileOptions) error {
		t.FileEnding = fileEnding
		return nil
	}
}

// TemporaryFile creates an ephemeral file. It returns the path to the created file as well as a cleanup function
// to remove the file later on. It is highly recommended to call this function once the fiel is no longer needed.
//
// The file ending can be modified using the WithFileEnding option.
//
// See also: TemporaryDirectory
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
	defer func() {
		_ = file.Close()
	}()

	cleanup := func() {
		_ = os.Remove(file.Name())
	}

	return file.Name(), cleanup, nil
}
