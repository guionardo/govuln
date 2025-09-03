package commands

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/guionardo/govuln/internal/params"
	"github.com/stretchr/testify/assert"
)

func TestGetRoot(t *testing.T) {
	params.CURRENT_PATH, _ = filepath.Abs("../..")
	root := GetRoot()
	assert.NotNil(t, root)
}

func Test_outputValidator(t *testing.T) {

	tests := []struct {
		name    string
		v       string
		wantErr bool
	}{
		{"valid_markdown", "markdown", false},
		{"valid_color", "color", false},
		{"valid_table", "table", false},
		{"invalid", "otherwise", true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := outputValidator(tt.v); (err != nil) != tt.wantErr {
				t.Errorf("outputValidator() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestRun(t *testing.T) {
	t.Run("run_help_should_succeed", func(t *testing.T) {
		os.Args = []string{os.Args[0], "--help"}
		assert.Equal(t, 0, Run())
	})
	t.Run("run_bad_argument_should_fail", func(t *testing.T) {
		os.Args = []string{os.Args[0], "--error"}
		assert.Equal(t, 1, Run())
	})
}
