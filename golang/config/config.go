package config

import (
	"io/ioutil"
	"path/filepath"
	"sync"

	"github.com/BurntSushi/toml"
)

var (
	builder     *Builder
	once        sync.Once
	defaultPath = "config.toml"
)

type Config interface{}

type Builder struct {
	Path string
}

var build = func() *Builder {
	once.Do(func() {
		builder = &Builder{
			Path: defaultPath,
		}
	})
	return builder
}()

func Path(path string) {
	build.Path = filepath.Clean(path)
}

func Load(config Config) error {
	b, err := ioutil.ReadFile(build.Path)
	if err != nil {
		return err
	}
	_, err = toml.Decode(string(b), config)
	return err
}
