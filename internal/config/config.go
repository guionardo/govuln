package config

import (
	"fmt"
	"os"
	"path"
	"path/filepath"

	"github.com/jedib0t/go-pretty/v6/table"
)

type Config struct {
	UserHomeDir       string
	StoreDefaultPath  string
	GoVulnCheckBinary string
	CurrentPath       string
	UseMarkDown       bool
}

const (
	AppName          = "FBM Fiscal GoVulnCheck"
	LocalStoreFolder = ".govulncheck"
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
		fmt.Printf("missing govulncheck binary!\nexpected location: %s\n\nrun this command to install: go install golang.org/x/vuln/cmd/govulncheck@latest\n", vcb)
		os.Exit(1)
	}

	currentPath, err := os.Getwd()
	if err != nil {
		fmt.Printf("failed to get current path: %v", err)
		os.Exit(1)
	}

	config = &Config{
		UserHomeDir:       userHomeDir,
		StoreDefaultPath:  filepath.Join(userHomeDir, LocalStoreFolder),
		GoVulnCheckBinary: vcb,
		CurrentPath:       currentPath,
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
