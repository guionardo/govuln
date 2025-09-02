package git

import (
	"fmt"
	"strings"
)

// ParseGitURL extracts host, owner, and repository components from a Git module path.
//
// This function parses Go module paths that reference Git repositories,
// typically in the format: "host.com/owner/repository" or "host.com/owner/repository/subpath".
// It's designed to work with internal Git hosting platforms and GitHub-style URLs.
//
// The parsing logic expects at least 3 path components separated by forward slashes:
//   - parts[0]: Git hosting platform hostname (e.g., "github.com", "gitlab.company.com")
//   - parts[1]: Repository owner/organization name
//   - parts[2]: Repository name (additional path components are ignored)
//
// Parameters:
//   - url: Go module path string (e.g., "github.com/melisource/my-repo")
//
// Returns:
//   - host: Git hosting platform hostname
//   - owner: Repository owner/organization
//   - repository: Repository name
//   - err: Non-nil if URL format is invalid (less than 3 components)
//
// Example:
//
//	host, owner, repo, err := ParseGitURL("github.com/melisource/govulncheck")
//	// Returns: "github.com", "melisource", "govulncheck", nil
//
// Note: This is a simplified parser and may need adjustments for complex URL formats
// or non-standard Git hosting configurations.
func ParseGitURL(url string) (host, owner, repository string, err error) {
	// Parse the git URL and extract the components
	// This is a simplified example and may need to be adjusted for different URL formats
	parts := strings.Split(url, "/")
	if len(parts) < 3 {
		err = fmt.Errorf("invalid git URL: %s", url)
		return
	}
	host = parts[0]
	owner = parts[1]
	repository = parts[2]
	return host, owner, repository, nil
}
