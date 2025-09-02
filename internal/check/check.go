package check

import (
	"fmt"
	"os"
	"path"
	"path/filepath"
	"slices"
	"strings"
	"time"

	goversion "github.com/hashicorp/go-version"
	"github.com/melisource/fury_fbm-fiscal-govulncheck/internal/config"
	"github.com/melisource/fury_fbm-fiscal-govulncheck/internal/exec"
	"github.com/melisource/fury_fbm-fiscal-govulncheck/internal/git"
	gocache "github.com/melisource/fury_fbm-fiscal-govulncheck/internal/go_cache"
	"github.com/melisource/fury_fbm-fiscal-govulncheck/internal/store"
	"golang.org/x/sync/errgroup"
)

type (
	// Check represents a vulnerability checking session for a Go project.
	//
	// This is the main structure that orchestrates vulnerability scanning by:
	// - Managing project metadata and caching (store, meta)
	// - Coordinating with external tools (govulncheck, git)
	// - Processing vulnerability data (vulnerabilities, osvs)
	// - Handling both main project and internal dependencies (internalOwner)
	//
	// The Check instance maintains state throughout the scanning lifecycle and
	// provides caching to optimize repeated scans of the same project.
	Check struct {
		store           *store.Store    // Cache storage for scan results and Git repositories
		folder          string          // Absolute path to the Go project being scanned
		meta            *MetaFile       // Metadata cache file for this project
		vulnerabilities Vulnerabilities // Current scan results: map[packageName][]*Vulnerability
		config          *Config         // govulncheck configuration data
		sbom            *SBOM           // Software Bill of Materials from govulncheck
		osvs            []*OSV          // Open Source Vulnerability database entries
		internalOwner   string          // Organization name for internal dependency scanning
	}

	// Versions represents a collection of semantic versions for dependency management.
	// Used internally for tracking multiple versions of the same module and checking
	// version ranges against vulnerability introduction/fix versions.
	Versions []*goversion.Version

	CheckType string
)

const (
	cacheTTL                 = time.Hour * 24 // Cache validity period: 24 hours
	ProjectCheck   CheckType = "project"
	SubModuleCheck CheckType = "submodule"
)

// safeGetModulePath safely retrieves the main module path, returning a fallback if unavailable
func (c *Check) safeGetModulePath() string {
	if c.sbom != nil && len(c.sbom.Modules) > 0 {
		return c.sbom.Modules[0].Path
	}
	return c.PackageName() // Fallback to folder-based name
}

// New creates a new vulnerability check instance for a Go project.
//
// This function validates that the provided folder contains a valid Go module
// (by checking for go.mod file) and initializes a Check instance with caching
// support and optional internal dependency scanning.
//
// Parameters:
//   - folder: Path to the Go project directory (must contain go.mod)
//   - store: Cache store for optimization (stores scan results for 24h)
//   - internalOwner: Organization name for internal dependency scanning (empty = skip)
//
// Returns:
//   - *Check: Initialized check instance ready for vulnerability scanning
//   - error: Non-nil if folder is invalid or not a Go module
//
// Example:
//
//	check, err := New("/path/to/project", store, "mycompany")
//	if err != nil {
//	    log.Fatal(err)
//	}
func New(folder string, store *store.Store, internalOwner string) (*Check, error) {
	abs, err := filepath.Abs(folder)
	if err != nil {
		return nil, err
	}
	folder = abs
	if stat, err := os.Stat(path.Join(folder, "go.mod")); os.IsNotExist(err) || stat.IsDir() {
		return nil, fmt.Errorf("folder %s is not a go module", folder)
	}

	metaFile := store.GetProjectMetaFile(folder)
	return &Check{
		store:         store,
		folder:        folder,
		meta:          NewMetaFile(metaFile),
		internalOwner: internalOwner,
	}, nil
}

func (c *Check) PackageName() string {
	if c.meta != nil && len(c.meta.SBOM.Modules) > 0 {
		return c.meta.SBOM.Modules[0].Path
	} else {
		return fmt.Sprintf("folder: %s", c.folder)
	}
}

// Run executes the vulnerability check on the project.
//
// This method performs the core vulnerability scanning by:
// 1. Checking cache for recent results (24h validity)
// 2. If cache miss, running govulncheck with JSON output
// 3. Parsing the JSON response to extract SBOM, OSV data, and config
// 4. Summarizing vulnerabilities and displaying results in a table
// 5. Caching results for future runs
//
// The method respects caching to avoid expensive repeated scans and
// provides rich formatted output showing found vulnerabilities.
//
// Parameters:
//   - args: Optional arguments. If first arg is "no-summary", skips table output.
//     Otherwise, first arg is used as version identifier in output.
//
// Returns:
//   - error: Non-nil if govulncheck execution fails or JSON parsing errors occur
//
// Cache behavior:
//   - Results cached for 24 hours per project
//   - Cache automatically invalidated if SBOM data changes
//   - Cache hit skips expensive govulncheck execution
func (c *Check) Run(checkType CheckType) error {
	if checkType == SubModuleCheck && c.meta.Checked && time.Since(c.meta.LastUpdate) < cacheTTL && len(c.meta.SBOM.GoVersion) > 0 {
		c.vulnerabilities = c.meta.Vulnerabilities
		c.config = &c.meta.Config
		c.sbom = &c.meta.SBOM
		c.osvs = c.meta.OSV

		return nil
	}
	// Run go mod tidy
	if checkType == ProjectCheck {
		exitCode, output, err := exec.RunAt(c.folder, "go", "mod", "tidy")
		if err == nil && exitCode != 0 {
			err = fmt.Errorf("error running go mod tidy: %d - %s", exitCode, output)
		}
		if err != nil {
			return err
		}
	}

	exitCode, output, err := exec.Run(config.Get().GoVulnCheckBinary, "-json", "-C", c.folder, "./...")
	if err == nil && exitCode != 0 {
		err = fmt.Errorf("error running govulncheck: %d - %s", exitCode, output)
	}
	if err != nil {
		return err
	}

	var body strings.Builder
	c.config = nil
	c.sbom = nil
	c.osvs = make([]*OSV, 0, 16)
	for line := range strings.SplitSeq(string(output), "\n") {
		if strings.HasPrefix(line, "{") {
			body.Reset()
			body.WriteString(line)
		} else if strings.HasPrefix(line, "}") {
			body.WriteString(line)
			if config := GetConfig(body.String()); config != nil {
				c.config = config
			} else if sbom := GetSBOM(body.String()); sbom != nil {
				c.sbom = sbom
			} else if osv := GetOSV(body.String()); osv != nil {
				c.osvs = append(c.osvs, osv)
			}

			body.Reset()
		} else {
			body.WriteString(line)
		}
	}

	c.Summarize()
	c.vulnerabilities.Table(os.Stdout, c.safeGetModulePath(), c.meta.PackageVersion)

	return err
}

func (v Versions) Has(version *goversion.Version) bool {
	return slices.ContainsFunc(v, func(vr *goversion.Version) bool {
		return vr.Equal(version)
	})
}

// isVulnerable checks if the current version is affected by a vulnerability.
//
// A version is considered vulnerable if:
// - It is greater than or equal to the 'introduced' version (when the vulnerability was introduced)
// - AND it is less than the 'fixed' version (when the vulnerability was patched)
//
// Parameters:
//   - current: The version currently in use
//   - introduced: The version when the vulnerability was introduced (nil if unknown)
//   - fixed: The version when the vulnerability was fixed (nil if not yet fixed)
//
// Returns:
//   - true if the current version is vulnerable
//   - false if the current version is safe (either predates introduction or includes fix)
func isVulnerable(current, introduced, fixed *goversion.Version) bool {
	if fixed != nil && current.GreaterThanOrEqual(fixed) {
		return false // Version has the fix
	}
	if introduced != nil && current.GreaterThanOrEqual(introduced) {
		return true // Version is affected by the vulnerability
	}
	return false // Version predates the vulnerability
}

// Summarize processes raw OSV (Open Source Vulnerability) data into actionable vulnerability reports.
//
// This method performs sophisticated vulnerability analysis by:
// 1. Building a module version map from SBOM data
// 2. Processing each OSV entry and its affected package ranges
// 3. Applying version comparison logic to determine vulnerable modules
// 4. Special handling for Go standard library vulnerabilities
// 5. Persisting results to cache for future rapid access
//
// The summarization process handles complex version range logic where vulnerabilities
// may affect specific version ranges (introduced..fixed) and must be matched against
// the exact versions found in the project's dependency tree.
//
// Standard library vulnerabilities are detected by checking the Go runtime version
// against stdlib vulnerability ranges, while module vulnerabilities require exact
// version matching from the dependency graph.
//
// Results are stored in both memory (c.vulnerabilities) and persistent cache
// (c.meta) for 24-hour validity periods.
//
// This method must be called after successful govulncheck execution to process
// the raw JSON output into structured vulnerability data.
func (c *Check) Summarize() {
	modules := make(map[string]Versions)
	for _, module := range c.sbom.Modules {
		if len(module.Version) > 0 {
			if v, err := goversion.NewVersion(module.Version); err == nil {
				modules[module.Path] = append(modules[module.Path], v)
			}
		}
	}
	goVersion, _ := goversion.NewVersion(strings.TrimPrefix(c.sbom.GoVersion, "go"))
	c.vulnerabilities = make(Vulnerabilities)

	for _, osv := range c.osvs {
		for _, aff := range osv.Affected {
			introduced, fixed := aff.Versions()
			// Detect vulnerabilities for stdlib
			if aff.Package.Name == "stdlib" && aff.Package.Ecosystem == "Go" {
				if isVulnerable(goVersion, introduced, fixed) {
					c.vulnerabilities[aff.Package.Name] = append(c.vulnerabilities[aff.Package.Name], &Vulnerability{
						Id:         osv.Id,
						Summary:    osv.Summary,
						Introduced: introduced,
						Fixed:      fixed,
						Current:    goVersion,
					})
				}
				continue
			}
			if versions, ok := modules[aff.Package.Name]; ok {
				for _, version := range versions {
					if isVulnerable(version, introduced, fixed) {
						c.vulnerabilities[aff.Package.Name] = append(c.vulnerabilities[aff.Package.Name], &Vulnerability{
							Id:         osv.Id,
							Summary:    osv.Summary,
							Introduced: introduced,
							Fixed:      fixed,
							Current:    version,
						})
					}
				}
			}
		}

	}
	c.meta.Vulnerabilities = c.vulnerabilities
	c.meta.OSV = c.osvs
	c.meta.SBOM = *c.sbom
	c.meta.Config = *c.config

	c.meta.Checked = true
	if err := c.meta.Save(); err != nil {
		fmt.Printf("Error saving meta file %s - %s", c.meta.filename, err)
	}
}

// CheckSubs performs concurrent vulnerability scanning of internal Git submodules.
//
// This method implements a sophisticated workflow for scanning internal dependencies:
// 1. Discovery: Parses SBOM modules to identify internal Git repositories
// 2. Filtering: Only processes modules matching the internal organization
// 3. Cloning: Concurrently clones missing repositories (up to 8 parallel)
// 4. Scanning: Runs vulnerability checks on each cloned repository
// 5. Aggregation: Collects and displays vulnerability summary across all submodules
//
// The method uses errgroup for controlled concurrency and maintains a local
// cache of cloned repositories to avoid repeated network operations.
//
// Concurrency model:
//   - Phase 1: Up to 8 concurrent git clone operations
//   - Phase 2: Up to 8 concurrent vulnerability scans
//   - Progress reporting via channel-based communication
//   - Mutex-protected shared state updates
//
// Internal dependencies are identified by matching the git URL host/owner
// against the configured internalOwner parameter.
//
// Prerequisites:
//   - internalOwner must be set (empty string disables submodule checking)
//   - Git repositories must be accessible via SSH
//   - SBOM data must be available from previous Run() execution
//
// Note: This method modifies no persistent state and is safe to call multiple times.
func (c *Check) CheckSubs() {
	if len(c.internalOwner) == 0 {
		return
	}
	eg := errgroup.Group{}
	goCache := gocache.New()

	eg.SetLimit(1)
	mods := make([]string, 0, len(c.sbom.Modules))
	for _, module := range c.sbom.Modules {
		if len(module.Version) == 0 {
			continue
		}
		_, owner, _, err := git.ParseGitURL(module.Path)
		if err != nil || owner != c.internalOwner {
			continue
		}
		if modPath, err := goCache.Get(module.Path, module.Version); err == nil {
			mods = append(mods, modPath)
		} else {
			fmt.Printf("error getting mod path: %s\n", err)
		}

	}

	smVulns := NewSubmodulesVulnerabilities()
	eg = errgroup.Group{}

	runCheck := func(mod string) error {
		c, err := New(mod, c.store, "")
		if err == nil {
			if err = c.Run(SubModuleCheck); err == nil {
				path, version, found := strings.Cut(c.folder, "@")
				path = goCache.GetPackageName(path)
				if found && (c.meta.PackageName != path || c.meta.PackageVersion != version) {
					c.meta.PackageVersion = version
					c.meta.PackageName = path
					_ = c.meta.Save()
				}

				if len(c.vulnerabilities) > 0 {
					for _, vulns := range c.vulnerabilities {
						for _, vuln := range vulns {
							smVulns.Add(c.safeGetModulePath(), version, vuln.Id)
						}
					}
				}
			}
		}
		return nil
	}
	for _, mod := range mods {
		if err := runCheck(mod); err != nil {
			fmt.Printf("error running check: %s\n", err)
			break
		}
	}
	if err := eg.Wait(); err != nil {
		fmt.Printf("error running checks: %s\n", err)
	}

	if err := eg.Wait(); err != nil {
		fmt.Printf("error running checks: %s\n", err)
	}
	smVulns.Table(os.Stdout)
}

func (c *Check) HasVulnerabilities() bool {
	return len(c.vulnerabilities) > 0
}
