package store

import (
	"fmt"
	"io/fs"
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"

	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/melisource/fury_fbm-fiscal-govulncheck/internal/config"
	"github.com/melisource/fury_fbm-fiscal-govulncheck/internal/tools/pathtools"
	"gopkg.in/yaml.v3"
)

type (
	Store struct {
		path string
	}
	MetaFile struct {
		PackageName        string    `yaml:"package_name"`
		PackageVersion     string    `yaml:"package_version"`
		LastUpdate         time.Time `yaml:"last_update"`
		HasVulnerabilities bool      `yaml:"has_vulnerabilities"`
	}
)

func New(path string) (*Store, error) {
	var err error
	if len(path) == 0 {
		path = config.Get().StoreDefaultPath
	} else {
		path, err = filepath.Abs(path)
		if err != nil {
			return nil, err
		}
	}
	if err = pathtools.CreatePath(filepath.Join(path, "projects")); err != nil {
		return nil, err
	}

	return &Store{
		path: path,
	}, nil
}

func (s *Store) Path(pieces ...string) (string, error) {
	p := filepath.Join(append([]string{s.path}, pieces...)...)
	return p, pathtools.CreatePath(p)
}

func (s *Store) ShowInfo() {
	fmt.Printf("Store path: %s\n", s.path)
	var fileCount, dirCount, fileSize int64
	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	t.SetStyle(table.StyleColoredBlackOnBlueWhite)
	t.SetTitle("Store informations")
	t.AppendHeader(table.Row{"Path", s.path})
	t.AppendHeader(table.Row{"Type", "Name", "Version", "Last Update", "Vulnerabilities"})

	err := filepath.WalkDir(s.path, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			dirCount++
		} else {
			if metaFile := ReadMetaFile(path); metaFile != nil {
				w := strings.Split(metaFile.PackageName, "melisource/")
				module := w[len(w)-1]
				if len(metaFile.PackageVersion) > 0 {
					t.AppendRow(table.Row{"Package", module, metaFile.PackageVersion, metaFile.LastUpdate.Format(time.DateTime), metaFile.HasVulnerabilities})
				} else {
					t.AppendRow(table.Row{"Project", module, "", metaFile.LastUpdate.Format(time.DateTime), metaFile.HasVulnerabilities})
				}
			}
			fileCount++
			if fi, err := d.Info(); err == nil {
				fileSize += fi.Size()
			}
		}
		return nil
	})
	t.AppendFooter(table.Row{"Files", fileCount, "Directories", dirCount, "Size", fileSize})
	if err != nil {
		t.AppendFooter(table.Row{"Error", err.Error()})
	}
	config.Render(t)

}

func (s *Store) Clear() {
	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	t.SetStyle(table.StyleColoredBlackOnBlueWhite)
	t.SetTitle("Store cleaning")
	t.AppendHeader(table.Row{"Directory", "Status"})
	dirEntries, err := os.ReadDir(s.path)
	if err != nil {
		t.AppendRow(table.Row{"ERROR", err.Error()})
		config.Render(t)
		os.Exit(1)
	}
	for _, dirEntry := range dirEntries {
		if dirEntry.IsDir() {
			err = os.RemoveAll(path.Join(s.path, dirEntry.Name()))
			if err != nil {
				t.AppendRow(table.Row{dirEntry.Name(), err.Error()})

			} else {
				t.AppendRow(table.Row{dirEntry.Name(), "Removed"})
			}
		}
	}
	config.Render(t)
}

func (s *Store) GetProjectMetaFile(projectPath string) string {
	normalized := strings.ReplaceAll(projectPath, "/", "_") + ".yaml"
	return path.Join(s.path, "projects", normalized)
}

func ReadMetaFile(filePath string) *MetaFile {
	if !strings.HasSuffix(filePath, ".yaml") {
		return nil
	}
	var meta MetaFile
	if f, err := os.Open(filePath); err == nil {
		err = yaml.NewDecoder(f).Decode(&meta)
		_ = f.Close()
		if err != nil {
			return nil
		}
	}

	if len(meta.PackageName) == 0 {
		return nil
	}
	return &meta
}
