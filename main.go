// Package main implements GoVuln, a comprehensive vulnerability scanner for Go projects.
//
// This tool wraps Go's official govulncheck utility with enhanced features including:
// - Intelligent caching (24-hour validity) to optimize repeated scans
// - Concurrent scanning of internal Git submodules and dependencies
// - Rich table output with color coding and Markdown support
// - Pre-commit hook integration for automated security checking
// - Organizational-specific dependency analysis
//
// The tool is designed for enterprise environments where Go projects depend on
// multiple internal repositories and require systematic vulnerability tracking
// across the entire dependency graph.
//
// Usage:
//
//	govuln [flags]
//
// Key features:
//
//	-just-warn: Warning mode (don't fail on vulnerabilities)
//	-dont-check-subs: Skip internal dependency scanning
//	-path: Custom project path to scan
//	-internal-owner: Organization name for dependency filtering
//
// The application requires govulncheck to be installed and accessible at
// ~/go/bin/govulncheck. Install with:
//
//	go install golang.org/x/vuln/cmd/govulncheck@latest
package main

import (
	"os"

	"github.com/guionardo/govuln/internal/commands"
)

func main() {
	os.Exit(commands.Run())
}
