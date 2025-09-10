package store

import (
	"os"
	"path"
	"reflect"
	"strings"
	"testing"

	_ "embed"

	"github.com/guionardo/govuln/internal/params"
	"github.com/stretchr/testify/assert"
)

//go:embed govuln.yaml
var govulnYaml []byte

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

	metaFileName := s.GetProjectMetaFile("github.com/guionardo/govuln")
	m := &MetaFile{
		Filename:       metaFileName,
		PackageName:    "github.com/guionardo/govuln",
		PackageVersion: "v0.0.0",
	}
	assert.NoError(t, s.SaveMetaFile(m))

	m2 := s.ReadMetaFile("github.com/guionardo/govuln")
	assert.NotNil(t, m2)
	assert.Equal(t, m.PackageName, m2.PackageName)

	assert.Error(t, s.SaveMetaFile(nil))

}

func TestStore_ReadMetaFile(t *testing.T) {
	type fields struct {
		path          string
		internalOwner string
	}
	type args struct {
		projectPath string
	}
	tests := []struct {
		name         string
		fields       fields
		args         args
		wantMetaFile *MetaFile
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Store{
				path:          tt.fields.path,
				internalOwner: tt.fields.internalOwner,
			}
			if gotMetaFile := s.ReadMetaFile(tt.args.projectPath); !reflect.DeepEqual(gotMetaFile, tt.wantMetaFile) {
				t.Errorf("Store.ReadMetaFile() = %v, want %v", gotMetaFile, tt.wantMetaFile)
			}
		})
	}
}

func TestReadMetaFile(t *testing.T) {
	t.Run("invalid_filename_should_fail", func(t *testing.T) {
		got := ReadMetaFile("meta.bad")
		assert.Nil(t, got)
	})
	t.Run("create_and_read_should_succeed", func(t *testing.T) {
		tmp := t.TempDir()
		m := &MetaFile{
			Filename:       path.Join(tmp, "metafile.yaml"),
			PackageName:    "github.com/guionardo/govuln",
			PackageVersion: "v0.0.0",
		}
		assert.NoError(t, m.Save())
		m2 := ReadMetaFile(m.Filename)
		assert.Equal(t, m.PackageName, m2.PackageName)
	})
	t.Run("invalid_content_should_fail", func(t *testing.T) {
		tmp := t.TempDir()
		metaFile := path.Join(tmp, "meta.yaml")
		assert.NoError(t, os.WriteFile(metaFile, []byte(`---
test: null
"`), 0644))
		m := ReadMetaFile(metaFile)
		assert.Nil(t, m)
	})

	t.Run("empty_package_name_should_fail", func(t *testing.T) {
		tmp := t.TempDir()
		m := &MetaFile{
			Filename: path.Join(tmp, "metafile.yaml"),
		}
		assert.NoError(t, m.Save())
		m2 := ReadMetaFile(m.Filename)
		assert.Nil(t, m2)
	})
}

func TestShowInfoAndClear(t *testing.T) {
	tmp := t.TempDir()
	s, err := New(tmp, "guionardo")
	if !assert.NotNil(t, s) || !assert.NoError(t, err) {
		return
	}
	metafile := s.GetProjectMetaFile("github.com/guionardo/govuln")
	if !assert.NoError(t, os.WriteFile(metafile, govulnYaml, 0644)) {
		return
	}

	assert.NotPanics(t, func() {
		s.ShowInfo()
	})
	assert.NotPanics(t, func() {
		s.Clear()
	})
}
