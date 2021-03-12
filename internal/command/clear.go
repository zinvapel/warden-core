package command

import (
	"github.com/spf13/cobra"
	"github.com/zinvapel/warden-core/pkg/environment"
)

func Clean() *cobra.Command {
	return &cobra.Command{
		Use:   "clean",
		Short: "Clean all jobs",
		Run: func(cmd *cobra.Command, args []string) {
			if env := environment.Unwrap(cmd); env != nil {
				env.Clear()
			} else {
				cmd.PrintErrln("Unable to get environment")
			}
		},
	}
}
