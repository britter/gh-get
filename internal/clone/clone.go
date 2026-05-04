package clone

import (
	"os"

	"github.com/britter/gh-get/internal/github"
	"github.com/cli/go-gh/v2/pkg/auth"
	git "github.com/go-git/go-git/v5"
	gitconfig "github.com/go-git/go-git/v5/config"
	githttp "github.com/go-git/go-git/v5/plumbing/transport/http"
)

// Clone clones target.Repository into clonePath. If target.Fork is set, the
// cloned origin remote is renamed to upstream and the fork is added as origin.
func Clone(target github.CloneTarget, clonePath string) error {
	cloneURL := "https://github.com/" + target.Repository.Owner + "/" + target.Repository.Name + ".git"

	opts := &git.CloneOptions{
		URL:      cloneURL,
		Progress: os.Stderr,
	}

	token, _ := auth.TokenForHost("github.com")
	if token != "" {
		opts.Auth = &githttp.BasicAuth{
			Username: "x-access-token",
			Password: token,
		}
	}

	r, err := git.PlainClone(clonePath, false, opts)
	if err != nil {
		return err
	}

	if target.Fork != nil {
		return setupForkRemotes(r, target.Fork)
	}
	return nil
}

func setupForkRemotes(r *git.Repository, fork *github.Repository) error {
	origin, err := r.Remote("origin")
	if err != nil {
		return err
	}
	upstreamURLs := origin.Config().URLs

	if err := r.DeleteRemote("origin"); err != nil {
		return err
	}
	if _, err := r.CreateRemote(&gitconfig.RemoteConfig{Name: "upstream", URLs: upstreamURLs}); err != nil {
		return err
	}
	forkURL := "https://github.com/" + fork.Owner + "/" + fork.Name + ".git"
	if _, err := r.CreateRemote(&gitconfig.RemoteConfig{Name: "origin", URLs: []string{forkURL}}); err != nil {
		return err
	}
	return nil
}
