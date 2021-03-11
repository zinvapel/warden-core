package warden

import (
	"errors"
	"github.com/spf13/cobra"
	"github.com/zinvapel/warden-core/internal"
	"github.com/zinvapel/warden-core/internal/command"
	"github.com/zinvapel/warden-core/internal/command/play"
	"github.com/zinvapel/warden-core/pkg/environment"
	"os"
)

const (
	AppName = "warden"
	DefaultHome = "/.warden"
)

func NewWarden() *internal.Warden {
	warden := &internal.Warden{Cobra: &cobra.Command{Use: AppName}}

	warden.Register(command.Doctor())
	warden.Register(command.Nip())
	warden.Register(command.Play())

	return warden
}

func DefaultEnvironment() (*environment.Environment, error) {
	if home, exist := os.LookupEnv("HOME"); exist {
		env := &environment.Environment{
			Home: home + DefaultHome,
			ConfigPaths: []string{home + DefaultHome},
			Cache: true,
		}

		env.Init()

		env.Walkers["script"] = &play.Script{}

		return env, nil
	}

	return nil, errors.New("$HOME does not exist")
}