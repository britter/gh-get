package github

import (
	"fmt"
	"strings"
)

type Repository struct {
	Owner string
	Name  string
}

func Parse(definition string) (Repository, error) {
	definition = strings.TrimPrefix(definition, "https://github.com/")
	definition = strings.TrimPrefix(definition, "git@github.com:")
	definition = strings.TrimSuffix(definition, ".git")

	parts := strings.Split(definition, "/")
	if len(parts) != 2 {
		return Repository{}, fmt.Errorf("Invalid repository definition: %s", definition)
	}

	owner, name := parts[0], parts[1]
	if owner == "" || name == "" {
		return Repository{}, fmt.Errorf("Invalid repository definition: %s", definition)
	}

	return Repository{parts[0], parts[1]}, nil
}
