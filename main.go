package main

import (
	"errors"
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/britter/gh-get/internal/github"
	"github.com/cli/go-gh/v2/pkg/auth"
	"github.com/cli/go-gh/v2/pkg/prompter"
	git "github.com/go-git/go-git/v5"
	githttp "github.com/go-git/go-git/v5/plumbing/transport/http"
)

func main() {
	forkFlag := flag.Bool("fork", false, "Fork the repository before cloning (--fork=false to skip prompt and clone original)")
	flag.Parse()

	var fork *bool
	flag.Visit(func(f *flag.Flag) {
		if f.Name == "fork" {
			fork = forkFlag
		}
	})

	ghFolder := getGhFolder()
	repositoryDefinition, err := getRepository()
	if err != nil {
		log.Fatal(err)
		return
	}

	repository, err := github.Parse(repositoryDefinition)
	if err != nil {
		log.Fatal(err)
		return
	}

	client, err := github.NewClient()
	if err != nil {
		log.Fatal(err)
		return
	}

	p := prompter.New(os.Stdin, os.Stdout, os.Stderr)
	target, err := github.ResolveCloneTarget(repository.Owner, repository.Name, fork, client, p)
	if err != nil {
		if errors.Is(err, github.ErrCancelled) {
			return
		}
		log.Fatal(err)
		return
	}

	cloneURL := "https://github.com/" + target.Owner + "/" + target.Name + ".git"
	clonePath := ghFolder + "/" + target.Owner + "/" + target.Name

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

	if _, err = git.PlainClone(clonePath, false, opts); err != nil {
		log.Fatal(err)
	}
}

func getRepository() (string, error) {
	args := flag.Args()
	if len(args) < 1 {
		return "", fmt.Errorf("No repository given for cloning. Please specify the repository to clone in '<owner>/<repository>' format.")
	}
	if len(args) > 1 {
		return "", fmt.Errorf("To many arguments. Please specify a single repository to clone in '<owner>/<repository>' format.")
	}
	return args[0], nil
}

func getGhFolder() string {
	repostoriesFolder := getenv("GH_GET_FOLDER", "github")
	fallbackRoot := os.Getenv("HOME") + "/" + repostoriesFolder
	return getenv("GH_GET_ROOT", fallbackRoot)
}

func getenv(key, fallback string) string {
	value := os.Getenv(key)
	if len(value) == 0 {
		return fallback
	}
	return value
}
