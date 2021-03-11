package registry

import "github.com/zinvapel/warden-core/pkg/environment"

type Source struct {
	Name        string
	Type        string
	Url         string
	Token       string
	ConfigRepos []Repo `mapstructure:"repos"`
}

type Repo struct {
	Name string
	Group string
	Url string
}

type Client interface {
	GetRepos(source Source) []Repo
}

type Walker interface {
	Walk(source *Source, extra string, env *environment.Environment) error
}