package store

import (
	"path"
	"strings"
	"testing"

	"github.com/guionardo/govuln/internal/params"
	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	t.Run("empty_path_get_default", func(t *testing.T) {
		s, err := New("", "")
		assert.NotNil(t, s)
		assert.NoError(t, err)
	})

	s, err := New(path.Join(t.TempDir(), params.LocalStoreFolder), "guionardo")
	assert.NoError(t, err)
	assert.True(t, strings.HasSuffix(s.path, params.LocalStoreFolder))

	p, err := s.Path("test", "folder")
	assert.NoError(t, err)
	assert.True(t, strings.HasSuffix(p, path.Join(params.LocalStoreFolder, "test", "folder")))
}
