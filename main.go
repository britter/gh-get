package main

import (
	"fmt"
	"log"
	"os"

	"github.com/cli/go-gh"
)

func main() {
	repositoriesRoot := getRepositoriesRoot()
	repository, err := getRepository()
	if err != nil {
		log.Fatal(err)
		return
	}

	repoClone := []string{"repo", "clone", repository, repositoriesRoot + "/" + repository}
	stdOut, stdErr, err := gh.Exec(repoClone...)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(stdOut.String())
	fmt.Println(stdErr.String())
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

func getRepositoriesRoot() string {
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
