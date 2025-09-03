package params

import (
	"os"
	"path"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_getGoVulnCheckTool(t *testing.T) {
	if !assert.NoError(t, readGoEnvs(), "expected go envs to be available") {
		return
	}

	// Try to remove current govulncheck tool
	gvc := path.Join(GOBIN, "govulncheck")
	if stat, err := os.Stat(gvc); err == nil && !stat.IsDir() {
		if err := os.Rename(gvc, gvc+".bkp"); err == nil {
			defer func() {
				if _, err := os.Stat(gvc); os.IsNotExist(err) {
					if err := os.Rename(gvc+".bkp", gvc); err != nil {
						t.Logf("WARNING: govulncheck need to be reinstalled!\nrun go install golang.org/x/vuln/cmd/govulncheck@latest")
					}
				} else {
					_ = os.Remove(gvc + ".bkp")
				}
			}()
		}
	}
	t.Run("simulates_install_error", func(t *testing.T) {
		GOVULNCHECK_VERSION = "BUG"
		err := getGoVulnCheckTool()
		assert.Error(t, err)
	})
	t.Run("install_suceeded", func(t *testing.T) {
		GOVULNCHECK_VERSION = "latest"
		err := getGoVulnCheckTool()
		assert.NoError(t, err)
	})

	GOVULNCHECK_VERSION = "latest"
}
