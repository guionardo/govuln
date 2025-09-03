package params

import (
	"fmt"
	"os/exec"
	"path"
	"strings"
)

var (
	GOMODCACHE string // /home/$USER/go/pkg/mod
	GOPATH     string // /home/$USER/go
	GOBIN      string // /home/$USER/go/bin
)

func readGoEnvs() error {
	cmd := exec.Command("go", "env")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return err
	}
	for line := range strings.SplitSeq(string(output), "\n") {
		words := strings.SplitN(strings.TrimSpace(line), "=", 2)
		if len(words) == 2 {
			value, _ := strings.CutPrefix(words[1], "'")
			value, _ = strings.CutSuffix(value, "'")
			switch words[0] {
			case "GOMODCACHE":
				GOMODCACHE = value
			case "GOPATH":
				GOPATH = value
			case "GOBIN":
				GOBIN = value
			}
		}
	}
	if len(GOPATH) == 0 {
		return fmt.Errorf("empty value for GOPATH")
	}
	if len(GOBIN) == 0 {
		GOBIN = path.Join(GOPATH, "bin")
	}

	return nil
}
