package internal

import (
	ctx "context"
	"github.com/spf13/cobra"
	"github.com/zinvapel/warden-core/pkg/environment"
)

// Prefer internal
type Warden struct {
	Env *environment.Environment
	Cobra *cobra.Command
}

func (w *Warden) WithEnv(env *environment.Environment) *Warden {
	w.Env = env
	return w
}

func (w *Warden) Run() error {
	return w.Cobra.ExecuteContext(ctx.WithValue(ctx.Background(), "environment", w.Env))
}

func (w *Warden) Register(command *cobra.Command) {
	w.Cobra.AddCommand(command)
}
