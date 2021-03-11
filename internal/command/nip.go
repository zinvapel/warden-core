package command

import (
	"errors"
	"fmt"
	"github.com/spf13/cobra"
	"net/url"
	"github.com/zinvapel/warden-core/pkg/environment"
	"github.com/zinvapel/warden-core/pkg/registry"
	"os"
)

func Nip() *cobra.Command {
	nip := &cobra.Command{
		Use:   "nip {type} {name} {url} {token} [-g GROUP_NAME]",
		Short: "Set up repositories resource",
		Args: func(cmd *cobra.Command, args []string) error {
			if env := environment.Unwrap(cmd); env != nil {
				if err := cobra.ExactArgs(4)(cmd, args); err != nil {
					return err
				}

				if _, ok := env.Registries[args[0]]; !ok {
					return errors.New(fmt.Sprintf("should be in %v", env.Registries))
				}

				if _, err := url.ParseRequestURI(args[2]); err != nil {
					return err
				}

				return nil
			} else {
				return errors.New("unable to get environment")
			}
		},
		Run: func(cmd *cobra.Command, args []string) {
			if env := environment.Unwrap(cmd); env != nil {
				src := registry.Source{
					Type:  args[0],
					Name:  args[1],
					Url:   args[2],
					Token: args[3],
				}

				fmt.Println("Obtaining repos")
				for _, singleRepo := range env.Registries[src.Type].GetRepos(src) {
					if groupFlag := cmd.Flag("group"); groupFlag != nil && groupFlag.Value.String() != "" {
						if singleRepo.Group != groupFlag.Value.String() {
							continue
						}
					}

					fmt.Printf("+ %s\n", singleRepo.Name)
					src.ConfigRepos = append(src.ConfigRepos, singleRepo)
				}

				fmt.Printf("Saving source with %d repos\n", len(src.ConfigRepos))
				env.Config.Set("sources." + src.Name, src)

				if err := env.Config.WriteConfig(); err != nil {
					cmd.PrintErrf("Unable to write config '%s'\n", err)
					os.Exit(1)
				}

				if err := env.SetUp(); err != nil {
					cmd.PrintErrf("Unable to doctor '%s'\n", err)
					os.Exit(1)
				}
			} else {
				cmd.PrintErrln("Unable to get environment")
			}
		},
	}

	nip.Flags().StringP("group", "g", "", "Filter repo by group name")

	return nip
}