package gocache

import (
	"fmt"
	"iter"
	"os"
	"path"
	"strings"

	"github.com/guionardo/govuln/internal/params"
)

type GoCache struct {
	modPath string
}

func New() *GoCache {
	return &GoCache{
		modPath: params.GOMODCACHE,
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
func (g *GoCache) GetPackages(internalOwner string) iter.Seq[string] {
	return func(yield func(string) bool) {
		// first level, domains
		domains, err := os.ReadDir(g.modPath)
		if err != nil {
			fmt.Printf("WARNING: error when reading data from go cache %s : %v", g.modPath, err)
			return
		}
		for _, domain := range domains {
			if domain.IsDir() && strings.Contains(domain.Name(), ".") {
				owners, _ := os.ReadDir(path.Join(g.modPath, domain.Name()))
			ownersFor:
				for owner := range owners {
					if owners[owner].Name() == internalOwner {
						repos, _ := os.ReadDir(path.Join(g.modPath, domain.Name(), internalOwner))
						for repo := range repos {
							if !yield(path.Join(domain.Name(), internalOwner, repos[repo].Name())) {
								return
							}
						}
						break ownersFor
					}
				}
			}
		}

	}
}
