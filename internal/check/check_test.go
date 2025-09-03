package check

import (
	"path/filepath"
	"testing"

	"github.com/guionardo/govuln/internal/params"
	"github.com/guionardo/govuln/internal/store"
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
