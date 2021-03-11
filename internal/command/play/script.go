package play

import (
	"encoding/json"
	"errors"
	"github.com/zinvapel/warden-core/pkg/environment"
	"github.com/zinvapel/warden-core/pkg/registry"
	"github.com/zinvapel/warden-core/pkg/utils/fs"
	"os"
	"sort"
)

type Script struct {}

func (s *Script) Walk(source *registry.Source, env *environment.Environment, extra string) error {
	settings := make(map[string]string)

	if err := json.Unmarshal([]byte(extra), &settings); err != nil {
		return err
	}

	if _, ok := settings["script"]; !ok {
		return errors.New("script not found")
	}

	if dirPath, err := env.WarmJobFor(source); err != nil {
		return err
	} else {
		os.Setenv("WARDEN_SOURCE_NAME", source.Name)
		os.Setenv("WARDEN_SOURCE_URL", source.Url)
		os.Setenv("WARDEN_SOURCE_TOKEN", source.Token)
		os.Setenv("WARDEN_DIR", dirPath)

		repos := make([]registry.Repo, 0)

		if filter, ok := settings["filter"]; ok {
			for _, repo := range source.ConfigRepos {
				if output, err := fs.ExecIn(dirPath, filter); err != nil {
					return err
				} else {
					if output == "pass" {
						repos = append(repos, repo)
					}
				}
			}
		} else {
			repos = append(repos, source.ConfigRepos...)
		}

		if filter, ok := settings["sort"]; ok {
			forSort := &repoSort{
				Repos: repos,
				LessFunc: func(l, r registry.Repo) bool {
					if out, err := fs.ExecIn(dirPath, filter); err == nil {
						return false
					} else {
						return out == "1"
					}
				},
			}

			sort.Sort(forSort)
			repos = forSort.Repos
		}

		for _, repo := range repos {
			os.Setenv("WARDEN_REPO_NAME", repo.Name)
			os.Setenv("WARDEN_REPO_URL", repo.Url)

			if err := fs.ContinueIn(dirPath, settings["script"]); err != nil {
				return err
			}
		}
	}

	return nil
}

type repoSort struct {
	Repos []registry.Repo
	LessFunc func(l, r registry.Repo) bool
}

func (s *repoSort) Len() int {
	return len(s.Repos)
}

func (s *repoSort) Less(i, j int) bool {
	return s.LessFunc(s.Repos[i], s.Repos[j])
}

func (s *repoSort) Swap(i, j int) {
	s.Repos[i], s.Repos[j] = s.Repos[j], s.Repos[i]
}
