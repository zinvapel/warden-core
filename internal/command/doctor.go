package command

import (
	"github.com/spf13/cobra"
	"github.com/zinvapel/warden-core/pkg/environment"
)

func Doctor() *cobra.Command {
	return &cobra.Command{
		Use:   "doctor",
		Short: "Verify and set up warden environment",
		Run: func(cmd *cobra.Command, args []string) {
			if env := environment.Unwrap(cmd); env != nil {
				if err := env.SetUp(); err != nil {
					cmd.PrintErrf("Unable to doctor '%s'\n", err)
				}
			} else {
				cmd.PrintErrln("Unable to get environment")
			}
		},
	}
}
