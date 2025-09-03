package check

import (
	"path/filepath"
	"testing"

	"github.com/guionardo/govuln/internal/params"
	"github.com/guionardo/govuln/internal/store"
	"github.com/hashicorp/go-version"
	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	store, _ := store.New("", "guionardo")
	params.CURRENT_PATH, _ = filepath.Abs("../..")
	c, err := New("../..", store, "guionardo")
	if !assert.NoError(t, err) || !assert.NotNil(t, c) {
		return
	}
	err = c.Run(ProjectCheck)
	assert.NoError(t, err)
	c.CheckSubs()
}

func TestVersions_Has(t *testing.T) {
	v1, _ := version.NewVersion("v1.0.0")
	v2, _ := version.NewVersion("v1.2.0")

	versions := Versions{v2}
	assert.True(t, versions.Has(*v2))
	assert.False(t, versions.Has(*v1))
}

func Test_isVulnerable(t *testing.T) {
	introduced, _ := version.NewVersion("v0.0.1")
	fixed, _ := version.NewVersion("v0.1.0")

	current, _ := version.NewVersion("v0.0.3")
	assert.True(t, isVulnerable(current, introduced, fixed))

	current, _ = version.NewVersion("v0.0.0")
	assert.False(t, isVulnerable(current, introduced, fixed))
}
