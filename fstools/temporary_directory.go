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

func WithPrefix(prefix string) function.Option[temporaryDirectoryOptions] {
	return func(options *temporaryDirectoryOptions) error {
		options.Prefix = prefix
		return nil
	}
}

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
