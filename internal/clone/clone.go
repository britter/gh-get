package clone

import (
	"fmt"
	"io"
	"os"
	"slices"
	"strings"
	"time"

	"github.com/britter/gh-get/internal/github"
	"github.com/cli/go-gh/v2/pkg/auth"
	git "github.com/go-git/go-git/v5"
	gitconfig "github.com/go-git/go-git/v5/config"
	"github.com/go-git/go-git/v5/plumbing/transport"
	githttp "github.com/go-git/go-git/v5/plumbing/transport/http"
)

// Clone clones target.Repository into clonePath. If target.Fork is set, the
// cloned origin remote is renamed to upstream and the fork is added as origin.
// diag receives verbose diagnostic messages; pass io.Discard to suppress them.
func Clone(target github.CloneTarget, clonePath string, diag io.Writer) error {
	cloneURL := "https://github.com/" + target.Repository.Owner + "/" + target.Repository.Name + ".git"

	opts := &git.CloneOptions{
		URL:  cloneURL,
		Auth: authMethod(diag),
	}

	fmt.Fprintf(diag, "Cloning %s into %s\n", cloneURL, clonePath)

	var r *git.Repository
	var err error

	if diag != io.Discard {
		opts.Progress = os.Stderr
		r, err = git.PlainClone(clonePath, false, opts)
	} else {
		fmt.Fprint(os.Stderr, "Cloning...")
		stop := make(chan struct{})
		go func() {
			ticker := time.NewTicker(3 * time.Second)
			defer ticker.Stop()
			for {
				select {
				case <-ticker.C:
					fmt.Fprint(os.Stderr, ".")
				case <-stop:
					return
				}
			}
		}()
		r, err = git.PlainClone(clonePath, false, opts)
		close(stop)
		fmt.Fprintln(os.Stderr, " done")
	}

	if err != nil {
		return err
	}

	fmt.Fprintln(os.Stdout, clonePath)

	if target.Fork != nil {
		fmt.Fprintf(diag, "Renaming origin to upstream, adding fork %s/%s as origin\n", target.Fork.Owner, target.Fork.Name)
		return setupForkRemotes(r, target.Fork, diag)
	}
	return nil
}

// Exists reports whether path is an existing git repository.
func Exists(path string) bool {
	_, err := git.PlainOpen(path)
	return err == nil
}

// Sync updates an existing clone at clonePath to match target and fast-forwards
// the current branch to upstream. If target.Fork is set and the clone currently
// points origin at the original, the remotes are reconfigured (origin ->
// upstream, fork -> origin). Local work is never clobbered: a diverged branch or
// dirty worktree is reported and skipped.
func Sync(target github.CloneTarget, clonePath string, p github.Prompter, diag io.Writer) error {
	r, err := git.PlainOpen(clonePath)
	if err != nil {
		return err
	}

	originalURL := "https://github.com/" + target.Repository.Owner + "/" + target.Repository.Name + ".git"

	if target.Fork != nil {
		origin, err := r.Remote("origin")
		if err != nil {
			return err
		}
		if clonesOriginal(origin, originalURL) {
			ok, err := p.Confirm(fmt.Sprintf("Reconfigure remotes (origin -> %s/%s, upstream -> %s/%s)?",
				target.Fork.Owner, target.Fork.Name, target.Repository.Owner, target.Repository.Name), true)
			if err != nil {
				return err
			}
			if !ok {
				fmt.Fprintln(os.Stderr, "Leaving remotes unchanged; fetching upstream only.")
			} else {
				fmt.Fprintf(diag, "Reconfiguring remotes for fork %s/%s\n", target.Fork.Owner, target.Fork.Name)
				if err := setupForkRemotes(r, target.Fork, diag); err != nil {
					return err
				}
			}
		}
	}

	// After reconfiguration the original lives in upstream; otherwise it's origin.
	remote := "origin"
	if _, err := r.Remote("upstream"); err == nil {
		remote = "upstream"
	}

	fmt.Fprintf(diag, "Synchronizing with %s\n", remote)
	w, err := r.Worktree()
	if err != nil {
		return err
	}
	err = w.Pull(&git.PullOptions{RemoteName: remote, Auth: authMethod(diag)})
	switch err {
	case nil, git.NoErrAlreadyUpToDate:
		// up to date or fast-forwarded
	case git.ErrNonFastForwardUpdate:
		fmt.Fprintln(os.Stderr, "Local branch has diverged from upstream, skipping fast-forward.")
	default:
		if strings.Contains(err.Error(), "worktree contains unstaged changes") {
			fmt.Fprintln(os.Stderr, "Working tree has local changes, skipping fast-forward.")
		} else {
			return err
		}
	}

	fmt.Fprintln(os.Stdout, clonePath)
	return nil
}

// clonesOriginal reports whether the origin remote points at the original repo.
func clonesOriginal(origin *git.Remote, originalURL string) bool {
	return slices.Contains(origin.Config().URLs, originalURL)
}

// authMethod returns GitHub token auth if a token is available, else nil.
func authMethod(diag io.Writer) transport.AuthMethod {
	token, _ := auth.TokenForHost("github.com")
	if token == "" {
		return nil
	}
	fmt.Fprintln(diag, "Using GitHub token for authentication")
	return &githttp.BasicAuth{Username: "x-access-token", Password: token}
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
