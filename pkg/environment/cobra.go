package environment

import (
	"github.com/spf13/cobra"
)

func Unwrap(cmd *cobra.Command) *Environment {
	if env, ok := cmd.Context().Value("environment").(*Environment); ok {
		return env
	}

	return nil
}
