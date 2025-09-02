package store

import (
	"path"
	"strings"
	"testing"

	"github.com/guionardo/govuln/internal/config"
	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	s, err := New(path.Join(t.TempDir(), config.LocalStoreFolder), "guionardo")
	assert.NoError(t, err)
	assert.True(t, strings.HasSuffix(s.path, config.LocalStoreFolder))

	p, err := s.Path("test", "folder")
	assert.NoError(t, err)
	assert.True(t, strings.HasSuffix(p, path.Join(config.LocalStoreFolder, "test", "folder")))
}
