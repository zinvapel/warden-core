package registry

type Source struct {
	Name        string
	Type        string
	Url         string
	Token       string
	ConfigRepos []Repo `yaml:"repos" mapstructure:"repos"`
}

type Repo struct {
	Id int
	Name string
	Group string
	Url string
}

type Client interface {
	GetRepos(source Source) []Repo
}