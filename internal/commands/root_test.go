package commands

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetRoot(t *testing.T) {
	root := GetRoot()
	assert.NotNil(t, root)
}
