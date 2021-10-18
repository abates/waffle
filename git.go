package waffle

import (
	"errors"
	"sort"

	"github.com/go-git/go-git/v5"
	gitconfig "github.com/go-git/go-git/v5/config"
	"github.com/go-git/go-git/v5/plumbing/object"
)

const (
	ErrNoGitRepo     = WaffleError("Git repo does not exist")
	ErrNoGitVersion  = WaffleError("Git repo does not have any semantic version tags")
	ErrInvalidGitURL = WaffleError("Git remote URL can't be parsed")
	ErrNoRepoName    = WaffleError("could not determine repo name")
	ErrOriginExists  = WaffleError("Remote 'origin' already exists")
)

func LoadGitMaintainer() (Maintainer, error) {
	m := Maintainer{}
	git, err := gitconfig.LoadConfig(gitconfig.GlobalScope)
	if err == nil {
		m.Name = git.User.Name
		m.Email = git.User.Email
	}
	return m, err
}

type GitRepo struct {
	repo *git.Repository
}

func OpenGit(dir string) (gr *GitRepo, err error) {
	repo, err := git.PlainOpen(dir)
	if err == nil {
		gr = &GitRepo{repo: repo}
	} else if errors.Is(err, git.ErrRepositoryNotExists) {
		err = ErrNoGitRepo
	}
	return
}

func InitGit(dir string) (gr *GitRepo, err error) {
	repo, err := git.PlainInit(dir, false)
	if err == nil {
		gr = &GitRepo{repo: repo}
	}
	return
}

func (gr *GitRepo) Remotes() (remotes map[string][]string, err error) {
	remotes = make(map[string][]string)
	r, err := gr.repo.Remotes()

	if err == nil {
		for _, remote := range r {
			c := remote.Config()
			for _, u := range c.URLs {
				remotes[c.Name] = append(remotes[c.Name], u)
			}
		}
	}
	return
}

func (gr *GitRepo) SetOrigin(url string) error {
	remotes, err := gr.Remotes()
	if err == nil {
		if _, found := remotes["origin"]; found {
			err = ErrOriginExists
		} else {
			_, err = gr.repo.CreateRemote(&gitconfig.RemoteConfig{
				Name: "origin",
				URLs: []string{"url"},
			})
		}
	}
	return err
}

func (gr *GitRepo) Versions() (versions []Version, err error) {
	tags, err := gr.repo.TagObjects()
	if err == nil {
		err = tags.ForEach(func(t *object.Tag) error {
			version := Version{}
			err = version.Set(t.Name)
			if err == nil {
				versions = append(versions, version)
			}
			return nil
		})
	}
	return
}

func (gr *GitRepo) CurrentVersion() (version Version, err error) {
	versions, err := gr.Versions()
	if err == nil {
		if len(versions) > 0 {
			sort.Sort(VersionList(versions))
			version = versions[len(versions)-1]
		} else {
			err = ErrNoGitVersion
		}
	}
	return
}
