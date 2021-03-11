package registry

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