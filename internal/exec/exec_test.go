package exec

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRun(t *testing.T) {
	gotExitCode, gotOutput, err := Run("ls", "-la")
	assert.NoError(t, err)
	assert.Equal(t, 0, gotExitCode)
	assert.NotEmpty(t, gotOutput)
}
