package clone

import (
	"io"
	"testing"

	"github.com/britter/gh-get/internal/github"
	git "github.com/go-git/go-git/v5"
	gitconfig "github.com/go-git/go-git/v5/config"
)

func TestExists(t *testing.T) {
	dir := t.TempDir()
	if Exists(dir) {
		t.Fatalf("empty dir should not be reported as a repo")
	}
	if _, err := git.PlainInit(dir, false); err != nil {
		t.Fatalf("init: %v", err)
	}
	if !Exists(dir) {
		t.Fatalf("initialized repo should be reported as existing")
	}
}

// When the existing clone points origin at the original and a fork is requested,
// Sync reconfigures remotes: origin -> fork, upstream -> original.
func TestSyncReconfiguresRemotesForFork(t *testing.T) {
	dir := t.TempDir()
	r, err := git.PlainInit(dir, false)
	if err != nil {
		t.Fatalf("init: %v", err)
	}
	originalURL := "https://github.com/owner/repo.git"
	if _, err := r.CreateRemote(&gitconfig.RemoteConfig{Name: "origin", URLs: []string{originalURL}}); err != nil {
		t.Fatalf("create origin: %v", err)
	}

	target := github.CloneTarget{
		Repository: github.Repository{Owner: "owner", Name: "repo"},
		Fork:       &github.Repository{Owner: "me", Name: "repo"},
	}
	// Pull will fail (no worktree/commits/network), but remote reconfiguration
	// happens first and is what we assert. Ignore the sync error here.
	_ = Sync(target, dir, io.Discard)

	upstream, err := r.Remote("upstream")
	if err != nil {
		t.Fatalf("expected upstream remote: %v", err)
	}
	if got := upstream.Config().URLs[0]; got != originalURL {
		t.Errorf("upstream URL = %q, want %q", got, originalURL)
	}
	origin, err := r.Remote("origin")
	if err != nil {
		t.Fatalf("expected origin remote: %v", err)
	}
	forkURL := "https://github.com/me/repo.git"
	if got := origin.Config().URLs[0]; got != forkURL {
		t.Errorf("origin URL = %q, want %q", got, forkURL)
	}
}

// When origin already points at a fork (not the original), remotes are left alone.
func TestSyncLeavesForkCloneRemotesAlone(t *testing.T) {
	dir := t.TempDir()
	r, err := git.PlainInit(dir, false)
	if err != nil {
		t.Fatalf("init: %v", err)
	}
	forkURL := "https://github.com/me/repo.git"
	if _, err := r.CreateRemote(&gitconfig.RemoteConfig{Name: "origin", URLs: []string{forkURL}}); err != nil {
		t.Fatalf("create origin: %v", err)
	}

	target := github.CloneTarget{
		Repository: github.Repository{Owner: "owner", Name: "repo"},
		Fork:       &github.Repository{Owner: "me", Name: "repo"},
	}
	_ = Sync(target, dir, io.Discard)

	if _, err := r.Remote("upstream"); err == nil {
		t.Errorf("upstream remote should not be created when origin already a fork")
	}
	origin, err := r.Remote("origin")
	if err != nil {
		t.Fatalf("origin remote missing: %v", err)
	}
	if got := origin.Config().URLs[0]; got != forkURL {
		t.Errorf("origin URL changed to %q, want %q", got, forkURL)
	}
}
