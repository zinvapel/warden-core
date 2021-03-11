package environment

import (
	"errors"
	"os"
)

const (
	DefaultHome = "/.warden"
)

func DefaultEnvironment() (*Environment, error) {
	if home, exist := os.LookupEnv("HOME"); exist {
		env := &Environment{
			Home: home + DefaultHome,
			ConfigPaths: []string{home + DefaultHome},
			Cache: true,
		}

		env.Init()

		return env, nil
	}

	return nil, errors.New("$HOME does not exist")
}
