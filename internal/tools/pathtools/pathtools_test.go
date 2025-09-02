package pathtools

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDirExists(t *testing.T) {
	tmp := t.TempDir()
	assert.True(t, DirExists(tmp))
	assert.False(t, DirExists(tmp+"/not-exists"))
}

func TestCreatePath(t *testing.T) {
	tmp := t.TempDir()
	assert.False(t, DirExists(tmp+"/not-exists"))
	assert.NoError(t, CreatePath(tmp+"/not-exists"))
	assert.True(t, DirExists(tmp+"/not-exists"))
}
