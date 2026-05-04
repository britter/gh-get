package clone

import (
	"fmt"
	"io"
	"os"

	"github.com/britter/gh-get/internal/github"
	"github.com/cli/go-gh/v2/pkg/auth"
	git "github.com/go-git/go-git/v5"
	gitconfig "github.com/go-git/go-git/v5/config"
	githttp "github.com/go-git/go-git/v5/plumbing/transport/http"
)

// Clone clones target.Repository into clonePath. If target.Fork is set, the
// cloned origin remote is renamed to upstream and the fork is added as origin.
// diag receives verbose diagnostic messages; pass io.Discard to suppress them.
func Clone(target github.CloneTarget, clonePath string, diag io.Writer) error {
	cloneURL := "https://github.com/" + target.Repository.Owner + "/" + target.Repository.Name + ".git"

	opts := &git.CloneOptions{
		URL:      cloneURL,
		Progress: os.Stderr,
	}

	token, _ := auth.TokenForHost("github.com")
	if token != "" {
		fmt.Fprintln(diag, "Using GitHub token for authentication")
		opts.Auth = &githttp.BasicAuth{
			Username: "x-access-token",
			Password: token,
		}
	}

	fmt.Fprintf(diag, "Cloning %s into %s\n", cloneURL, clonePath)
	r, err := git.PlainClone(clonePath, false, opts)
	if err != nil {
		return err
	}

	if target.Fork != nil {
		fmt.Fprintf(diag, "Renaming origin to upstream, adding fork %s/%s as origin\n", target.Fork.Owner, target.Fork.Name)
		return setupForkRemotes(r, target.Fork, diag)
	}
	return nil
}

func setupForkRemotes(r *git.Repository, fork *github.Repository, diag io.Writer) error {
	origin, err := r.Remote("origin")
	if err != nil {
		return err
	}
	upstreamURLs := origin.Config().URLs
	fmt.Fprintf(diag, "Creating upstream remote: %v\n", upstreamURLs)

	if err := r.DeleteRemote("origin"); err != nil {
		return err
	}
	if _, err := r.CreateRemote(&gitconfig.RemoteConfig{Name: "upstream", URLs: upstreamURLs}); err != nil {
		return err
	}
	forkURL := "https://github.com/" + fork.Owner + "/" + fork.Name + ".git"
	fmt.Fprintf(diag, "Creating origin remote: %s\n", forkURL)
	if _, err := r.CreateRemote(&gitconfig.RemoteConfig{Name: "origin", URLs: []string{forkURL}}); err != nil {
		return err
	}
	return nil
}
