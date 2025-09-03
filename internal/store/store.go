package store

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"

	pathtools "github.com/guionardo/go/pkg/path_tools"
	"github.com/guionardo/govuln/internal/config"
	"github.com/guionardo/govuln/internal/output"
	"github.com/jedib0t/go-pretty/v6/table"
	"gopkg.in/yaml.v3"
)

type (
	Store struct {
		path          string
		internalOwner string
	}
)

func New(path, internalOwner string) (store *Store, err error) {
	if len(path) == 0 {
		path = config.Get().StoreDefaultPath
	} else {
		if p, err := filepath.Abs(path); err == nil {
			path = p
		}
	}
	if err = pathtools.CreatePath(filepath.Join(path, "projects")); err == nil {
		store = &Store{
			path:          path,
			internalOwner: internalOwner,
		}
	}

	return store, err
}

func (s *Store) ReadMetaFile(projectPath string) (metaFile *MetaFile) {
	return ReadMetaFile(s.GetProjectMetaFile(projectPath))
}

func (s *Store) SaveMetaFile(metaFile *MetaFile) error {
	if metaFile != nil {
		return metaFile.Save()
	}
	return errors.New("null metafile - can not save")
}

func (s *Store) Path(pieces ...string) (string, error) {
	p := filepath.Join(append([]string{s.path}, pieces...)...)
	return p, pathtools.CreatePath(p)
}

func (s *Store) ShowInfo() {
	fmt.Printf("Store path: %s\n", s.path)
	var fileCount, dirCount, fileSize int64
	t := output.New(os.Stdout, "Store informations", false)
	t.AppendHeader(table.Row{"Path", s.path})
	t.AppendHeader(table.Row{"Type", "Name", "Version", "Last Update", "Vulnerabilities"})

	err := filepath.WalkDir(s.path, func(filePath string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			dirCount++
		} else {
			if metaFile := ReadMetaFile(filePath); metaFile != nil {
				if metaFile.HasVulnerabilities {
					t.SetWithError(true)
				}
				w := strings.Split(metaFile.PackageName, path.Join(s.internalOwner, ""))
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
	t.Render()

}

func (s *Store) Clear() {
	t := output.New(os.Stdout, "Store cleaning", false)
	t.AppendHeader(table.Row{"Directory", "Status"})
	dirEntries, err := os.ReadDir(s.path)
	if err != nil {
		t.AppendRow(table.Row{"ERROR", err.Error()})
		t.Render()
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
	t.Render()
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
