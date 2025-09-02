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
	"context"
	"log"
	"os"

	"github.com/guionardo/govuln/internal/commands"
)

func main() {
	cmd := commands.GetRoot()
	if err := cmd.Run(context.Background(), os.Args); err != nil {
		log.Fatal(err)
	}
}

// func mainOld() {

// 	var (
// 		justWarn      = flag.Bool("just-warn", false, "Just warn about vulnerabilities without failing")
// 		dontCheckSubs = flag.Bool("dont-check-subs", false, "Don't check submodules")
// 		path          = flag.String("path", config.Get().CurrentPath, "Path of the project to check vulnerabilities")
// 		storePath     = flag.String("store", config.Get().StoreDefaultPath, "Path of the store for caching checks of internal packages")
// 		storeInfo     = flag.Bool("store-info", false, "Show store information")
// 		storeClear    = flag.Bool("store-clear", false, "Clear store")
// 		internalOwner = flag.String("internal-owner", "guionardo", "Internal owner")
// 		useMarkdown   = flag.Bool("use-markdown", false, "Use Markdown for table output")
// 	)

// 	flag.Parse()

// 	config.Get().UseMarkDown = *useMarkdown

// 	store, err := store.New(*storePath, *internalOwner)
// 	if err != nil {
// 		fmt.Printf("%s - failed to create store\n\t[%s]\n", config.AppName, *storePath)
// 		os.Exit(1)
// 	}
// 	if *storeInfo {
// 		store.ShowInfo()
// 		os.Exit(0)
// 	}
// 	if *storeClear {
// 		store.Clear()
// 	}
// 	if !pathtools.DirExists(*path) {
// 		fmt.Printf("%s - path does not exist\n\t[%s]\n", config.AppName, *path)
// 		os.Exit(2)
// 	}
// 	if len(*internalOwner) == 0 {
// 		*dontCheckSubs = true
// 	} else if *dontCheckSubs {
// 		*internalOwner = ""
// 	}
// 	chk, err := check.New(*path, store, *internalOwner)

// 	if err != nil {
// 		fmt.Printf("%s - failed to create check\n\t[%s]\n", config.AppName, *path)
// 		os.Exit(1)
// 	}
// 	fmt.Printf("%s - %s - checking vulnerabilities", config.AppName, chk.PackageName())
// 	if !*dontCheckSubs {
// 		fmt.Print(" and submodules")
// 	}
// 	if *justWarn {
// 		fmt.Print(" [just-warn]")
// 	}
// 	fmt.Println()

// 	err = chk.Run(check.ProjectCheck)
// 	if err != nil {
// 		fmt.Printf("%s - failed to run check\n\t[%s]\n", config.AppName, *path)
// 		os.Exit(1)
// 	}
// 	if !*dontCheckSubs {
// 		chk.CheckSubs()
// 	}
// 	if chk.HasVulnerabilities() && !*justWarn {
// 		os.Exit(1)
// 	}

// }
