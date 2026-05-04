package github

import (
	"bytes"
	"errors"
	"io"
	"strings"
	"testing"
)

type fakeClient struct {
	info      RepoInfo
	infoErr   error
	forkOwner string
	forkName  string
	forkErr   error
}

func (f *fakeClient) RepoInfo(_, _ string) (RepoInfo, error) {
	return f.info, f.infoErr
}

func (f *fakeClient) Fork(_, _ string) (string, string, error) {
	return f.forkOwner, f.forkName, f.forkErr
}

type fakePrompter struct {
	answer bool
	err    error
}

func (f *fakePrompter) Confirm(_ string, _ bool) (bool, error) {
	return f.answer, f.err
}

func boolPtr(b bool) *bool { return &b }

func repoPtr(owner, name string) *Repository {
	r := Repository{owner, name}
	return &r
}

func TestResolveCloneTarget(t *testing.T) {
	tests := []struct {
		name           string
		fork           *bool
		client         *fakeClient
		prompter       *fakePrompter
		expectedTarget CloneTarget
		expectedErr    error
	}{
		{
			name: "has write access, no flag: clone original",
			fork: nil,
			client: &fakeClient{
				info: RepoInfo{AllowForking: true, HasPushAccess: true},
			},
			prompter:       &fakePrompter{},
			expectedTarget: CloneTarget{Repository: Repository{"owner", "repo"}},
		},
		{
			name: "no write access, forking allowed, user says yes: fork",
			fork: nil,
			client: &fakeClient{
				info:      RepoInfo{AllowForking: true, HasPushAccess: false},
				forkOwner: "me",
				forkName:  "repo",
			},
			prompter:       &fakePrompter{answer: true},
			expectedTarget: CloneTarget{Repository: Repository{"owner", "repo"}, Fork: repoPtr("me", "repo")},
		},
		{
			name: "no write access, forking allowed, user says yes, repo was renamed: use fork name",
			fork: nil,
			client: &fakeClient{
				info:      RepoInfo{AllowForking: true, HasPushAccess: false},
				forkOwner: "me",
				forkName:  "new-name",
			},
			prompter:       &fakePrompter{answer: true},
			expectedTarget: CloneTarget{Repository: Repository{"owner", "repo"}, Fork: repoPtr("me", "new-name")},
		},
		{
			name: "no write access, forking allowed, user says no: cancelled",
			fork: nil,
			client: &fakeClient{
				info: RepoInfo{AllowForking: true, HasPushAccess: false},
			},
			prompter:    &fakePrompter{answer: false},
			expectedErr: ErrCancelled,
		},
		{
			name: "no write access, forking disabled: clone original with warning",
			fork: nil,
			client: &fakeClient{
				info: RepoInfo{AllowForking: false, HasPushAccess: false},
			},
			prompter:       &fakePrompter{},
			expectedTarget: CloneTarget{Repository: Repository{"owner", "repo"}},
		},
		{
			name: "--fork=true, forking allowed: fork",
			fork: boolPtr(true),
			client: &fakeClient{
				info:      RepoInfo{AllowForking: true, HasPushAccess: false},
				forkOwner: "me",
				forkName:  "repo",
			},
			prompter:       &fakePrompter{},
			expectedTarget: CloneTarget{Repository: Repository{"owner", "repo"}, Fork: repoPtr("me", "repo")},
		},
		{
			name: "--fork=true, forking allowed, has write access: fork anyway",
			fork: boolPtr(true),
			client: &fakeClient{
				info:      RepoInfo{AllowForking: true, HasPushAccess: true},
				forkOwner: "me",
				forkName:  "repo",
			},
			prompter:       &fakePrompter{},
			expectedTarget: CloneTarget{Repository: Repository{"owner", "repo"}, Fork: repoPtr("me", "repo")},
		},
		{
			name: "--fork=true, forking allowed, repo was renamed: use fork name",
			fork: boolPtr(true),
			client: &fakeClient{
				info:      RepoInfo{AllowForking: true, HasPushAccess: false},
				forkOwner: "me",
				forkName:  "new-name",
			},
			prompter:       &fakePrompter{},
			expectedTarget: CloneTarget{Repository: Repository{"owner", "repo"}, Fork: repoPtr("me", "new-name")},
		},
		{
			name: "--fork=false, has write access: clone original",
			fork: boolPtr(false),
			client: &fakeClient{
				info: RepoInfo{AllowForking: true, HasPushAccess: true},
			},
			prompter:       &fakePrompter{},
			expectedTarget: CloneTarget{Repository: Repository{"owner", "repo"}},
		},
		{
			name: "--fork=false, no write access: clone original without prompting",
			fork: boolPtr(false),
			client: &fakeClient{
				info: RepoInfo{AllowForking: true, HasPushAccess: false},
			},
			prompter:       &fakePrompter{},
			expectedTarget: CloneTarget{Repository: Repository{"owner", "repo"}},
		},
		{
			name: "no write access, forking allowed, prompter fails: propagate error",
			fork: nil,
			client: &fakeClient{
				info: RepoInfo{AllowForking: true, HasPushAccess: false},
			},
			prompter:    &fakePrompter{err: errors.New("prompt failed")},
			expectedErr: errors.New("prompt failed"),
		},
		{
			name: "--fork=true, forking disabled: error",
			fork: boolPtr(true),
			client: &fakeClient{
				info: RepoInfo{AllowForking: false, HasPushAccess: false},
			},
			prompter:    &fakePrompter{},
			expectedErr: errors.New("repository owner/repo does not allow forking"),
		},
		{
			name: "API error fetching repo info: propagate error",
			fork: nil,
			client: &fakeClient{
				infoErr: errors.New("API error"),
			},
			prompter:    &fakePrompter{},
			expectedErr: errors.New("API error"),
		},
		{
			name: "fork API call fails: propagate error",
			fork: boolPtr(true),
			client: &fakeClient{
				info:    RepoInfo{AllowForking: true},
				forkErr: errors.New("fork failed"),
			},
			prompter:    &fakePrompter{},
			expectedErr: errors.New("fork failed"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			target, err := ResolveCloneTarget("owner", "repo", tt.fork, tt.client, tt.prompter, io.Discard)

			if tt.expectedErr != nil {
				if err == nil {
					t.Fatalf("expected error %q, got nil", tt.expectedErr)
				}
				if tt.expectedErr == ErrCancelled {
					if !errors.Is(err, ErrCancelled) {
						t.Errorf("expected ErrCancelled, got %v", err)
					}
				} else if err.Error() != tt.expectedErr.Error() {
					t.Errorf("expected error %q, got %q", tt.expectedErr, err)
				}
				return
			}

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if target.Repository.Owner != tt.expectedTarget.Repository.Owner {
				t.Errorf("expected repository owner %q, got %q", tt.expectedTarget.Repository.Owner, target.Repository.Owner)
			}
			if target.Repository.Name != tt.expectedTarget.Repository.Name {
				t.Errorf("expected repository name %q, got %q", tt.expectedTarget.Repository.Name, target.Repository.Name)
			}
			if tt.expectedTarget.Fork == nil {
				if target.Fork != nil {
					t.Errorf("expected no fork, got %v", target.Fork)
				}
			} else {
				if target.Fork == nil {
					t.Fatalf("expected fork %v, got nil", tt.expectedTarget.Fork)
				}
				if target.Fork.Owner != tt.expectedTarget.Fork.Owner {
					t.Errorf("expected fork owner %q, got %q", tt.expectedTarget.Fork.Owner, target.Fork.Owner)
				}
				if target.Fork.Name != tt.expectedTarget.Fork.Name {
					t.Errorf("expected fork name %q, got %q", tt.expectedTarget.Fork.Name, target.Fork.Name)
				}
			}
		})
	}
}

func TestResolveCloneTargetDiagnosticOutput(t *testing.T) {
	tests := []struct {
		name            string
		fork            *bool
		client          *fakeClient
		prompter        *fakePrompter
		expectedInDiag  []string
	}{
		{
			name: "has write access: reports push access",
			fork: nil,
			client: &fakeClient{
				info: RepoInfo{AllowForking: true, HasPushAccess: true},
			},
			prompter: &fakePrompter{},
			expectedInDiag: []string{
				"Fetching repository info for owner/repo",
				"Have push access, cloning original",
			},
		},
		{
			name: "--fork=false: reports skip",
			fork: boolPtr(false),
			client: &fakeClient{
				info: RepoInfo{AllowForking: true, HasPushAccess: true},
			},
			prompter: &fakePrompter{},
			expectedInDiag: []string{
				"Fetching repository info for owner/repo",
				"Skipping fork (--fork=false)",
			},
		},
		{
			name: "fork created: reports allow forking and push access flags",
			fork: boolPtr(true),
			client: &fakeClient{
				info:      RepoInfo{AllowForking: true, HasPushAccess: false},
				forkOwner: "me",
				forkName:  "repo",
			},
			prompter: &fakePrompter{},
			expectedInDiag: []string{
				"Fetching repository info for owner/repo",
				"Allow forking: true, has push access: false",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer
			ResolveCloneTarget("owner", "repo", tt.fork, tt.client, tt.prompter, &buf) //nolint:errcheck
			output := buf.String()
			for _, want := range tt.expectedInDiag {
				if !strings.Contains(output, want) {
					t.Errorf("expected diag output to contain %q, got:\n%s", want, output)
				}
			}
		})
	}
}
