package github

import (
	"errors"
	"fmt"
	"io"
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

// CloneTarget describes what to clone and how to configure remotes.
// When Fork is non-nil, Repository is the upstream and Fork is the user's fork.
type CloneTarget struct {
	Repository Repository
	Fork       *Repository
}

// ResolveCloneTarget determines which repository to clone, forking if necessary.
//
// fork controls forking behaviour:
//   - nil  — default: prompt when the user lacks write access
//   - true — always fork (skip prompt); error if forking is disabled
//   - false — never fork; clone the original directly without prompting
//
// When a fork is created, CloneTarget.Fork is set; the caller should clone
// Repository as upstream and add Fork as origin.
//
// diag receives verbose diagnostic messages; pass io.Discard to suppress them.
//
// Returns ErrCancelled if the user declines to fork.
func ResolveCloneTarget(owner, name string, fork *bool, client Client, prompter Prompter, diag io.Writer) (CloneTarget, error) {
	fmt.Fprintf(diag, "Fetching repository info for %s/%s\n", owner, name)
	info, err := client.RepoInfo(owner, name)
	if err != nil {
		return CloneTarget{}, err
	}
	fmt.Fprintf(diag, "Allow forking: %v, has push access: %v\n", info.AllowForking, info.HasPushAccess)

	if fork != nil && *fork {
		if !info.AllowForking {
			return CloneTarget{}, fmt.Errorf("repository %s/%s does not allow forking", owner, name)
		}
		return forkRepo(owner, name, client)
	}

	if fork != nil && !*fork {
		fmt.Fprintf(diag, "Skipping fork (--fork=false)\n")
		return CloneTarget{Repository: Repository{owner, name}}, nil
	}

	if !info.HasPushAccess {
		if !info.AllowForking {
			fmt.Fprintf(os.Stderr, "You do not have write access to %s/%s and forking is disabled, cloning original.\n", owner, name)
			return CloneTarget{Repository: Repository{owner, name}}, nil
		}
		answer, err := prompter.Confirm(fmt.Sprintf("You don't have write access to %s/%s. Fork it?", owner, name), true)
		if err != nil {
			return CloneTarget{}, err
		}
		if !answer {
			return CloneTarget{}, ErrCancelled
		}
		return forkRepo(owner, name, client)
	}

	fmt.Fprintf(diag, "Have push access, cloning original\n")
	return CloneTarget{Repository: Repository{owner, name}}, nil
}

func forkRepo(owner, name string, client Client) (CloneTarget, error) {
	forkOwner, forkName, err := client.Fork(owner, name)
	if err != nil {
		return CloneTarget{}, err
	}
	fork := Repository{forkOwner, forkName}
	return CloneTarget{Repository: Repository{owner, name}, Fork: &fork}, nil
}
