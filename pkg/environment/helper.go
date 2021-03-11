package environment

import (
	"fmt"
	"github.com/zinvapel/warden-core/pkg/registry"
	"github.com/zinvapel/warden-core/pkg/utils/fs"
)

type Walker interface {
	Walk(source *registry.Source, env *Environment, extra string) error
}

type FetchStrategy interface {
	FetchSourceInto(src *registry.Source, dir string) error
}

type GitStrategy struct {
	MainBranch string
}

func (g *GitStrategy) FetchSourceInto(src *registry.Source, baseDir string) error {
	if _, err := fs.GetOrCreateDir(baseDir); err != nil {
		return err
	}

	for _, repo := range src.ConfigRepos {
		fmt.Println("Processing repo", repo.Name)

		if fs.Exist(baseDir + repo.Name) {
			if err := fs.ContinueIn(baseDir + repo.Name, "git checkout -- ."); err != nil {
				return err
			}

			if err := fs.ContinueIn(baseDir + repo.Name, "git checkout " + g.MainBranch); err != nil {
				return err
			}

			if err := fs.ContinueIn(baseDir + repo.Name, "git pull origin " + g.MainBranch); err != nil {
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