package check

import (
	"testing"

	"github.com/melisource/fury_fbm-fiscal-govulncheck/internal/store"
	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	store, _ := store.New("")
	c, err := New("../..", store, "melisource")
	if !assert.NoError(t, err) || !assert.NotNil(t, c) {
		return
	}
	err = c.Run(ProjectCheck)
	assert.NoError(t, err)
	c.CheckSubs()
}
