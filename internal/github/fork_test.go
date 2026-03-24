package github

import (
	"errors"
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

func TestResolveCloneTarget(t *testing.T) {
	tests := []struct {
		name         string
		forceFork    bool
		client       *fakeClient
		prompter     *fakePrompter
		expectedRepo Repository
		expectedErr  error
	}{
		{
			name:      "has write access, no flag: clone original",
			forceFork: false,
			client: &fakeClient{
				info: RepoInfo{AllowForking: true, HasPushAccess: true},
			},
			prompter:     &fakePrompter{},
			expectedRepo: Repository{"owner", "repo"},
		},
		{
			name:      "no write access, forking allowed, user says yes: fork",
			forceFork: false,
			client: &fakeClient{
				info:      RepoInfo{AllowForking: true, HasPushAccess: false},
				forkOwner: "me",
				forkName:  "repo",
			},
			prompter:     &fakePrompter{answer: true},
			expectedRepo: Repository{"me", "repo"},
		},
		{
			name:      "no write access, forking allowed, user says yes, repo was renamed: use fork name",
			forceFork: false,
			client: &fakeClient{
				info:      RepoInfo{AllowForking: true, HasPushAccess: false},
				forkOwner: "me",
				forkName:  "new-name",
			},
			prompter:     &fakePrompter{answer: true},
			expectedRepo: Repository{"me", "new-name"},
		},
		{
			name:      "no write access, forking allowed, user says no: cancelled",
			forceFork: false,
			client: &fakeClient{
				info: RepoInfo{AllowForking: true, HasPushAccess: false},
			},
			prompter:    &fakePrompter{answer: false},
			expectedErr: ErrCancelled,
		},
		{
			name:      "no write access, forking disabled: clone original with warning",
			forceFork: false,
			client: &fakeClient{
				info: RepoInfo{AllowForking: false, HasPushAccess: false},
			},
			prompter:     &fakePrompter{},
			expectedRepo: Repository{"owner", "repo"},
		},
		{
			name:      "--fork flag, forking allowed: fork",
			forceFork: true,
			client: &fakeClient{
				info:      RepoInfo{AllowForking: true, HasPushAccess: false},
				forkOwner: "me",
				forkName:  "repo",
			},
			prompter:     &fakePrompter{},
			expectedRepo: Repository{"me", "repo"},
		},
		{
			name:      "--fork flag, forking allowed, has write access: fork anyway",
			forceFork: true,
			client: &fakeClient{
				info:      RepoInfo{AllowForking: true, HasPushAccess: true},
				forkOwner: "me",
				forkName:  "repo",
			},
			prompter:     &fakePrompter{},
			expectedRepo: Repository{"me", "repo"},
		},
		{
			name:      "--fork flag, forking allowed, repo was renamed: use fork name",
			forceFork: true,
			client: &fakeClient{
				info:      RepoInfo{AllowForking: true, HasPushAccess: false},
				forkOwner: "me",
				forkName:  "new-name",
			},
			prompter:     &fakePrompter{},
			expectedRepo: Repository{"me", "new-name"},
		},
		{
			name:      "no write access, forking allowed, prompter fails: propagate error",
			forceFork: false,
			client: &fakeClient{
				info: RepoInfo{AllowForking: true, HasPushAccess: false},
			},
			prompter:    &fakePrompter{err: errors.New("prompt failed")},
			expectedErr: errors.New("prompt failed"),
		},
		{
			name:      "--fork flag, forking disabled: error",
			forceFork: true,
			client: &fakeClient{
				info: RepoInfo{AllowForking: false, HasPushAccess: false},
			},
			prompter:    &fakePrompter{},
			expectedErr: errors.New("repository owner/repo does not allow forking"),
		},
		{
			name:      "API error fetching repo info: propagate error",
			forceFork: false,
			client: &fakeClient{
				infoErr: errors.New("API error"),
			},
			prompter:    &fakePrompter{},
			expectedErr: errors.New("API error"),
		},
		{
			name:      "fork API call fails: propagate error",
			forceFork: true,
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
			repo, err := ResolveCloneTarget("owner", "repo", tt.forceFork, tt.client, tt.prompter)

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
			if repo.Owner != tt.expectedRepo.Owner {
				t.Errorf("expected owner %q, got %q", tt.expectedRepo.Owner, repo.Owner)
			}
			if repo.Name != tt.expectedRepo.Name {
				t.Errorf("expected name %q, got %q", tt.expectedRepo.Name, repo.Name)
			}
		})
	}
}
