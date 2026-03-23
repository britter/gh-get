package main

import (
	"fmt"
	"log"
	"os"

	"github.com/britter/gh-get/internal/github"
	"github.com/cli/go-gh/v2/pkg/auth"
	git "github.com/go-git/go-git/v5"
	githttp "github.com/go-git/go-git/v5/plumbing/transport/http"
)

func main() {
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

	cloneURL := "https://github.com/" + repository.Owner + "/" + repository.Name + ".git"
	clonePath := ghFolder + "/" + repository.Owner + "/" + repository.Name

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
	args := os.Args[1:]
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
