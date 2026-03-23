package github

import (
	"testing"
)

func TestParse(t *testing.T) {
	tests := []struct {
		input         string
		expectedOwner string
		expectedName  string
		expectError   bool
	}{
		// Valid cases
		{"owner/name", "owner", "name", false},
		{"https://github.com/owner/name", "owner", "name", false},
		{"https://github.com/owner/name.git", "owner", "name", false},
		{"git@github.com:owner/name.git", "owner", "name", false},
		{"https://github.com/owner/name/tree/main", "owner", "name", false},
		{"https://github.com/owner/name/tree/feature/my-feature", "owner", "name", false},
		{"https://github.com/owner/name/tree/v1.2.3", "owner", "name", false},
		{"https://github.com/owner/name/commit/abc1234", "owner", "name", false},
		{"https://github.com/owner/name/blob/main/README.md", "owner", "name", false},
		{"https://github.com/owner/name/pull/42", "owner", "name", false},
		{"https://github.com/owner/name/issues/42", "owner", "name", false},

		// Invalid cases
		{"https://github.com/invalidinput", "", "", true}, // Missing /
		{"invalidinput", "", "", true},                    // Missing /
		{"https://github.com/", "", "", true},             // Empty path
		{"owner/", "", "", true},                          // Missing name
		{"/name", "", "", true},                           // Missing owner
	}

	for _, test := range tests {
		repository, err := Parse(test.input)

		if test.expectError {
			if err == nil {
				t.Errorf("Expected error for input %q, but got none", test.input)
			}
		} else {
			if err != nil {
				t.Errorf("Did not expect error fro input %q, but got: %v", test.input, err)
			}
			if repository.Owner != test.expectedOwner {
				t.Errorf("For input %q, expected owner %q, but got %q", test.input, test.expectedOwner, repository.Owner)
			}
			if repository.Name != test.expectedName {
				t.Errorf("For input %q, expected name %q, but got %q", test.input, test.expectedName, repository.Name)
			}
		}
	}
}
