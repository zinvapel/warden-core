package environment

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/zinvapel/warden-core/pkg/registry"
	"github.com/zinvapel/warden-core/pkg/utils/fs"
	"os"
)

const (
	varDir = "/var"
	jobDir = varDir + "/jobs"
	gitDir = "/git"

	configName = "config"
	configType = "yaml"
)

type Environment struct {
	JobId string

	Home string
	HomeDir *os.FileInfo
	ConfigPaths []string
	Config *viper.Viper

	Extensions []*cobra.Command
	Registries map[string]registry.Client
	Walkers map[string]Walker

	Cache bool
}

type Walker interface {
	Walk(source *registry.Source, env *Environment, extra string) error
}

func (env *Environment) Init() {
	env.JobId = uuid.New().String()

	home, _ := fs.GetOrCreateDir(env.Home)
	env.HomeDir = &home

	env.Config = viper.New()
	env.Config.SetConfigName(configName)
	env.Config.SetConfigType(configType)
	for _, path := range env.ConfigPaths {
		env.Config.AddConfigPath(path)
	}
	_ = env.Config.ReadInConfig()

	env.Registries = make(map[string]registry.Client)
	env.Extensions = make([]*cobra.Command, 0)
	env.Walkers = make(map[string]Walker)
}

func (env *Environment) WarmJobFor(src *registry.Source) (string, error) {
	if _, err := fs.GetOrCreateDir(env.Home + jobDir); err != nil {
		return "", err
	}

	dir := env.Home + jobDir + "/" + env.JobId

	if env.Cache {
		err := fs.ContinueInArgs(
			env.Home,
			"cp",
			"-R",
			env.Home + gitDir + "/.",
			dir + "/.",
		)

		if err != nil {
			return "", err
		}
	} else {
		if err := env.FetchSourceInto(src, dir); err != nil {
			return "", err
		}
	}

	return dir, nil
}

func (env *Environment) SetUp() error {
	if _, err := fs.GetOrCreateDir(env.Home); err != nil {
		return err
	}

	if _, err := fs.GetOrCreateDir(env.Home + varDir); err != nil {
		return err
	}

	if _, err := fs.GetOrCreateDir(env.Home + jobDir); err != nil {
		return err
	}

	fs.Remove(env.Home + jobDir + "/*")

	if _, err := fs.GetOrCreateDir(env.Home + gitDir); err != nil {
		return err
	}

	for _, path := range env.ConfigPaths {
		if _, err := fs.GetOrCreateFile(path + "/" + configName + "." + configType); err != nil {
			return err
		}
	}

	var sources map[string]*registry.Source

	if err := env.Config.UnmarshalKey("sources", &sources); err != nil {
		return err
	}

	for name, src := range sources {
		baseDir := env.Home + gitDir + "/" + name + "/"

		if err := env.FetchSourceInto(src, baseDir); err != nil {
			return err
		}
	}

	return nil
}

func (env *Environment) FetchSourceInto(src *registry.Source, baseDir string) error {
	if _, err := fs.GetOrCreateDir(baseDir); err != nil {
		return err
	}

	for _, repo := range src.ConfigRepos {
		fmt.Println("Processing repo", repo.Name)

		if fs.Exist(baseDir + repo.Name) {
			// @todo add login/pass or ssh_keys
			if err := fs.ContinueIn(baseDir + repo.Name, "git checkout master"); err != nil {
				return err
			}

			if err := fs.ContinueIn(baseDir + repo.Name, "git pull origin master"); err != nil {
				return err
			}
		} else {
			if err := fs.ContinueIn(baseDir, "git clone " + repo.Url + " " + repo.Name); err != nil {
				return err
			}
		}
	}

	return nil
}