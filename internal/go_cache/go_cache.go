package gocache

import (
	"fmt"
	"os"
	"path"
	"strings"
)

type GoCache struct {
	modPath string
}

func New() *GoCache {
	goPath := os.Getenv("GOPATH")
	if goPath == "" {
		goPath = os.Getenv("HOME") + "/go"
	}
	modPath := path.Join(goPath, "pkg", "mod")
	return &GoCache{
		modPath: modPath,
	}
}

func (g *GoCache) Get(packagePath string, version string) (string, error) {
	modPath := path.Join(g.modPath, fmt.Sprintf("%s@%s", packagePath, version))
	_, err := os.Stat(modPath)
	os.IsNotExist(err)
	return modPath, err
}

func (g *GoCache) GetPackageName(packagePath string) string {
	if after, found := strings.CutPrefix(packagePath, g.modPath); found {
		packagePath, _ = strings.CutPrefix(after, "/")
	}
	return packagePath
}
