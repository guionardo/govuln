package gocache

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	_ = os.Setenv("GOPATH", "")
	cache := New()
	if !assert.NotNil(t, cache) {
		return
	}
	p, err := cache.Get("github.com/guionardo/govuln", "v0.0.0")
	assert.Error(t, err)
	assert.NotEmpty(t, p)

	pn := cache.GetPackageName(p)
	assert.Equal(t, "github.com/guionardo/govuln@v0.0.0", pn)
}

func TestGoCache_GetPackages(t *testing.T) {
	cache := New()
	packages := make([]string, 0)
	assert.NotPanics(t, func() {
		for meta := range cache.GetPackages("guionardo") {
			packages = append(packages, meta)
		}

	})

}
