package warden

import (
	"github.com/spf13/cobra"
	"github.com/zinvapel/warden-core/internal"
	"github.com/zinvapel/warden-core/internal/command"
)

const AppName = "warden"

func NewWarden() *internal.Warden {
	warden := &internal.Warden{Cobra: &cobra.Command{Use: AppName}}

	warden.Register(command.Doctor())
	warden.Register(command.Nip())
	warden.Register(command.Play())

	return warden
}
