package config

import (
	"errors"
	"fmt"
	"os"
	"path"
	"path/filepath"

	"github.com/guionardo/govuln/internal/exec"
	"github.com/guionardo/govuln/internal/git"
	"github.com/jedib0t/go-pretty/v6/table"
)

type Config struct {
	UserHomeDir       string
	StoreDefaultPath  string
	GoVulnCheckBinary string
	CurrentPath       string
	UseMarkDown       bool
	InternalOwner     string
}

const (
	AppName          = "Guionardo's GoVulnCheck"
	LocalStoreFolder = ".govuln"
)

var config *Config

func init() {
	userHomeDir, err := os.UserHomeDir()
	if err != nil {
		fmt.Printf("failed to get user home directory: %v", err)
		os.Exit(1)
	}

	// govulncheck installation
	vcb := path.Join(userHomeDir, "go", "bin", "govulncheck")
	if stat, err := os.Stat(vcb); os.IsNotExist(err) || stat.IsDir() {
		if err := installGoVulnCheck(); err != nil {
			fmt.Printf("❌ missing govulncheck binary!\nexpected location: %s\n\nrun this command to install: go install golang.org/x/vuln/cmd/govulncheck@latest\n", vcb)
			os.Exit(1)
		}
	}

	currentPath, err := os.Getwd()
	if err != nil {
		fmt.Printf("failed to get current path: %v", err)
		os.Exit(1)
	}

	_, owner, _, _ := git.ParseGitFromLocalFolder(currentPath)

	config = &Config{
		UserHomeDir:       userHomeDir,
		StoreDefaultPath:  filepath.Join(userHomeDir, LocalStoreFolder),
		GoVulnCheckBinary: vcb,
		CurrentPath:       currentPath,
		InternalOwner:     owner,
	}
}

func Get() *Config {
	return config
}

func Render(table table.Writer) {
	if config.UseMarkDown {
		table.RenderMarkdown()
	} else {
		table.Render()
	}
	fmt.Println()
}

func installGoVulnCheck() error {
	fmt.Print("Trying to install latest govulncheck...")
	exitCode, output, err := exec.Run("go", "install", "golang.org/x/vuln/cmd/govulncheck@latest")
	if exitCode == 0 && err == nil {
		fmt.Println(" ✅")
		return nil
	}

	fmt.Printf("Exit code: %d - %s", exitCode, output)
	return errors.New("failed installing govulncheck")
}
