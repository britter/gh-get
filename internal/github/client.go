package github

import (
	"fmt"
	"os"
	"time"

	"github.com/cli/go-gh/v2/pkg/api"
)

type restClient struct {
	client *api.RESTClient
}

// NewClient creates a Client backed by the GitHub REST API.
func NewClient() (Client, error) {
	c, err := api.DefaultRESTClient()
	if err != nil {
		return nil, err
	}
	return &restClient{client: c}, nil
}

type repoResponse struct {
	AllowForking bool `json:"allow_forking"`
	Permissions  struct {
		Push  bool `json:"push"`
		Admin bool `json:"admin"`
	} `json:"permissions"`
}

func (c *restClient) RepoInfo(owner, name string) (RepoInfo, error) {
	var resp repoResponse
	if err := c.client.Get(fmt.Sprintf("repos/%s/%s", owner, name), &resp); err != nil {
		return RepoInfo{}, err
	}
	return RepoInfo{
		AllowForking:  resp.AllowForking,
		HasPushAccess: resp.Permissions.Push || resp.Permissions.Admin,
	}, nil
}

func (c *restClient) Fork(owner, name string) (string, string, error) {
	fmt.Fprintf(os.Stderr, "Forking %s/%s...\n", owner, name)

	var forkResp struct {
		Owner struct {
			Login string `json:"login"`
		} `json:"owner"`
		Name string `json:"name"`
	}
	if err := c.client.Post(fmt.Sprintf("repos/%s/%s/forks", owner, name), nil, &forkResp); err != nil {
		return "", "", err
	}

	forkOwner := forkResp.Owner.Login
	forkName := forkResp.Name

	fmt.Fprintf(os.Stderr, "Waiting for fork to become available...")
	if err := c.waitForFork(forkOwner, forkName, 30*time.Second); err != nil {
		fmt.Fprintln(os.Stderr)
		return "", "", err
	}
	fmt.Fprintln(os.Stderr, " done")
	return forkOwner, forkName, nil
}

func (c *restClient) waitForFork(forkOwner, name string, timeout time.Duration) error {
	deadline := time.Now().Add(timeout)
	for time.Now().Before(deadline) {
		var resp struct {
			Size int `json:"size"`
		}
		err := c.client.Get(fmt.Sprintf("repos/%s/%s", forkOwner, name), &resp)
		if err == nil && resp.Size > 0 {
			return nil
		}
		fmt.Fprint(os.Stderr, ".")
		time.Sleep(2 * time.Second)
	}
	return fmt.Errorf("timed out waiting for fork to become available")
}
