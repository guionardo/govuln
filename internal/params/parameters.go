package params

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"

	"github.com/guionardo/govuln/internal/git"
)

const (
	AppName          = "Guionardo's GoVulnCheck"
	LocalStoreFolder = ".govuln"
)

var (
	USER_HOME_DIR        string
	OUTPUT_TYPE          string
	GO_VULN_CHECK_BINARY string
	STORE_DEFAULT_PATH   string
	CURRENT_PATH         string
	INTERNAL_OWNER       = sync.OnceValue(func() string {
		_, owner, _, err := git.ParseGitFromLocalFolder(CURRENT_PATH)
		if err != nil {
			panic(fmt.Errorf("error getting internal owner: %v", err))
		}
		return owner
	})
)

func init() {
	runUntilError(getUserHomeDir, readGoEnvs, getGoVulnCheckTool, getCurrentPath)
	STORE_DEFAULT_PATH = filepath.Join(USER_HOME_DIR, LocalStoreFolder)
}

func getUserHomeDir() error {
	uhd, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("error getting user home dir: %v", err)
	}
	USER_HOME_DIR = uhd
	return nil
}

// func getInternalOwner() error {
// 	_, owner, _, err := git.ParseGitFromLocalFolder(CURRENT_PATH)
// 	if err != nil {
// 		return fmt.Errorf("error getting internal owner: %v", err)
// 	}
// 	INTERNAL_OWNER = owner
// 	return nil
// }

func getCurrentPath() error {
	if len(CURRENT_PATH) == 0 {
		currentPath, err := os.Getwd()
		if err != nil {
			return fmt.Errorf("error getting current path: %v", err)

		}
		CURRENT_PATH = currentPath
	}
	return nil
}

func runUntilError(funcs ...func() error) {
	for _, fn := range funcs {
		if err := fn(); err != nil {
			fmt.Printf("❌ %v\n", err)
			os.Exit(1)
		}
	}
}
