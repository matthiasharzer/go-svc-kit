package fstools

import (
	"fmt"
	"math/rand"
	"os"

	"github.com/matthiasharzer/service-kit/function"
)

var randomPatternsPrefixes = []string{
	"jupiter",
	"titan",
	"europa",
	"mars",
	"platypus",
	"badger",
	"kraken",
	"yeti",
	"capybara",
	"enterprise",
	"serenity",
	"nostromo",
	"tardis",
	"voyager",
	"abyss",
	"limbo",
	"vortex",
	"oblivion",
	"gizmo",
	"blackhole",
	"narnia",
	"doohickey",
	"purgatory",
}

func getRandomPrefix() string {
	randomIndex := rand.Intn(len(randomPatternsPrefixes))
	return randomPatternsPrefixes[randomIndex]
}

type temporaryDirectoryOptions struct {
	Prefix string
}

// WithPrefix configures the prefix used when creating a temporary directory. This inflicts
// the name of the folder inside the os's temporary directory.
func WithPrefix(prefix string) function.Option[temporaryDirectoryOptions] {
	return func(options *temporaryDirectoryOptions) error {
		options.Prefix = prefix
		return nil
	}
}

// TemporaryDirectory creates an ephemeral directory. It returns the name of the directory as well as
// a cleanup function, which when called removes the temporary directory.
//
// The prefix of the created directory name can be modified using the WithPrefix option. If omitted, a random prefix
// will be used.
//
// See also: TemporaryFile
func TemporaryDirectory(opts ...function.Option[temporaryDirectoryOptions]) (string, func(), error) {
	options := temporaryDirectoryOptions{
		Prefix: getRandomPrefix(),
	}
	err := function.ApplyOptionsInto(opts, &options)
	if err != nil {
		return "", nil, err
	}

	dir, err := os.MkdirTemp("", fmt.Sprintf("%s-*", options.Prefix))
	if err != nil {
		return "", nil, err
	}

	cleanup := func() {
		_ = os.RemoveAll(dir)
	}

	return dir, cleanup, nil
}
