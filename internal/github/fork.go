package github

import (
	"errors"
	"fmt"
	"os"
)

// ErrCancelled is returned when the user declines to fork.
var ErrCancelled = errors.New("cancelled")

// Prompter can prompt the user for confirmation.
type Prompter interface {
	Confirm(prompt string, defaultValue bool) (bool, error)
}

// RepoInfo contains the fields from the GitHub API needed for fork decisions.
type RepoInfo struct {
	AllowForking  bool
	HasPushAccess bool
}

// Client provides the GitHub API operations needed for fork decisions.
type Client interface {
	RepoInfo(owner, name string) (RepoInfo, error)
	// Fork creates a fork and returns the fork's owner and name, which may
	// differ from the original if the repository was renamed.
	Fork(owner, name string) (forkOwner, forkName string, err error)
}

// ResolveCloneTarget determines which repository to clone, forking if necessary.
//
// fork controls forking behaviour:
//   - nil  — default: prompt when the user lacks write access
//   - true — always fork (skip prompt); error if forking is disabled
//   - false — never fork; clone the original directly without prompting
//
// Returns ErrCancelled if the user declines to fork.
func ResolveCloneTarget(owner, name string, fork *bool, client Client, prompter Prompter) (Repository, error) {
	info, err := client.RepoInfo(owner, name)
	if err != nil {
		return Repository{}, err
	}

	if fork != nil && *fork {
		if !info.AllowForking {
			return Repository{}, fmt.Errorf("repository %s/%s does not allow forking", owner, name)
		}
		return forkRepo(owner, name, client)
	}

	if fork != nil && !*fork {
		return Repository{owner, name}, nil
	}

	if !info.HasPushAccess {
		if !info.AllowForking {
			fmt.Fprintf(os.Stderr, "You do not have write access to %s/%s and forking is disabled, cloning original.\n", owner, name)
			return Repository{owner, name}, nil
		}
		answer, err := prompter.Confirm(fmt.Sprintf("You don't have write access to %s/%s. Fork it?", owner, name), false)
		if err != nil {
			return Repository{}, err
		}
		if !answer {
			return Repository{}, ErrCancelled
		}
		return forkRepo(owner, name, client)
	}

	return Repository{owner, name}, nil
}

func forkRepo(owner, name string, client Client) (Repository, error) {
	forkOwner, forkName, err := client.Fork(owner, name)
	if err != nil {
		return Repository{}, err
	}
	return Repository{forkOwner, forkName}, nil
}
