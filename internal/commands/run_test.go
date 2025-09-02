package commands

import (
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_run(t *testing.T) {
	cmd := runCommand()
	t.Run("invalid_project_path_must_fail", func(t *testing.T) {
		assert.Error(t, run(t.Context(), cmd))
	})

	projectPath, _ = filepath.Abs("../../.") // Valid project path
	t.Run("invalid_store_path_must_fail", func(t *testing.T) {
		storePath = "/x_invalid" // Invalid store path
		assert.Error(t, run(t.Context(), cmd))
	})

	t.Run("valid_project_and_store_must_succeed", func(t *testing.T) {
		storePath = t.TempDir()
		assert.NoError(t, run(t.Context(), cmd))
	})

}
