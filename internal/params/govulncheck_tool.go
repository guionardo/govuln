package params

import (
	"errors"
	"fmt"
	"os"
	"path"

	"github.com/guionardo/govuln/internal/exec"
)

var GOVULNCHECK_VERSION = "latest"

func getGoVulnCheckTool() error {
	vcb := path.Join(GOBIN, "govulncheck")
	if stat, err := os.Stat(vcb); os.IsNotExist(err) || stat.IsDir() {
		if err := installGoVulnCheck(GOVULNCHECK_VERSION); err != nil {
			return fmt.Errorf("missing govulncheck binary!\nexpected location: %s\n\nrun this command to install: go install golang.org/x/vuln/cmd/govulncheck@latest", vcb)
		}
	}
	GO_VULN_CHECK_BINARY = vcb
	return nil
}

// govulncheck installation

func installGoVulnCheck(version string) error {
	fmt.Print("Trying to install latest govulncheck...")
	exitCode, output, err := exec.Run("go", "install", "golang.org/x/vuln/cmd/govulncheck@"+version)
	if exitCode == 0 && err == nil {
		fmt.Println(" ✅")
		return nil
	}

	fmt.Printf("Exit code: %d - %s", exitCode, output)
	return errors.New("failed installing govulncheck")
}
