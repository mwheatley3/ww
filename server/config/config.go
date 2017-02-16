package config

import (
	"github.com/mwheatley3/ww/server/log"

	"github.com/mwheatley3/ww/server/config/config"

	"fmt"
	"os"
	"os/signal"
	"syscall"
)

const (
	// Delimiter separates object levels in the config
	Delimiter = "__"
)

// Load loads config from each of the providers
// and finally from the environment
func Load(conf interface{}, providers ...config.Provider) error {
	c := make(chan os.Signal)
	go func() {
		for {
			select {
			case <-c:
				//Not using a logger here, because there is a cyclic dependency in the
				//load call sites that can't be broken.
				fmt.Printf("Config: \n\t%v\n", config.String(conf))
			}
		}
	}()
	signal.Notify(c, syscall.SIGUSR1)

	providers = append(providers, config.FromEnv())
	return config.Load(conf, Delimiter, providers...)
}

// MustLoad runs load an dies on an error
func MustLoad(conf interface{}, providers ...config.Provider) {
	if err := Load(conf, providers...); err != nil {
		log.Fatalf("Config load error: %s", err)
	}
}

// FromFileWithOverride is a config.Provider that loads config
// from a file and then tries to load from an override file
func FromFileWithOverride(path string) config.Provider {
	f1 := config.FromFile(path)
	f2 := config.FromFile(path + ".override")

	return func() ([]string, error) {
		d, err := f1()

		if err != nil {
			return nil, err
		}

		d2, err := f2()

		// don't care about .override not existing
		if err == nil {
			d = append(d, d2...)
		}

		return d, nil
	}
}
