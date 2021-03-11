package command

import (
	"errors"
	"github.com/spf13/cobra"
	"os"
	"github.com/zinvapel/warden-core/pkg/environment"
	"github.com/zinvapel/warden-core/pkg/registry"
)

func Play() *cobra.Command {
	play := &cobra.Command{
		Use:   "play {name} {play}",
		Short: "Interact with source",
		Args: func(cmd *cobra.Command, args []string) error {
			if err := cobra.ExactArgs(2)(cmd, args); err != nil {
				return err
			}

			if env := environment.Unwrap(cmd); env != nil {
				if !env.Config.IsSet("sources." + args[0]) {
					return errors.New("source does not exist")
				}
				if _, ok := env.Walkers[args[1]]; !ok {
					return errors.New("play does not exist")
				}
			}

			return nil
		},
		Run: func(cmd *cobra.Command, args []string) {
			if env := environment.Unwrap(cmd); env != nil {
				src := &registry.Source{}

				if err := env.Config.UnmarshalKey("sources." + args[0], src); err != nil {
					cmd.PrintErrln("Source does not exist")
					os.Exit(1)
				}

				err := env.Walkers[args[1]].Walk(src, env, cmd.Flag("extra").Value.String())
				if err != nil {
					cmd.PrintErrln("Error on play", err)
				}
			} else {
				cmd.PrintErrln("Unable to get environment")
			}
		},
	}

	play.Flags().StringP("extra", "e", "", "Extra data")

	return play
}