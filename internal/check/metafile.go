package check

import (
	"os"
	"time"

	"gopkg.in/yaml.v3"
)

type MetaFile struct {
	filename           string
	PackageName        string          `yaml:"package_name"`
	PackageVersion     string          `yaml:"package_version"`
	Vulnerabilities    Vulnerabilities `yaml:"vulnerabilities"`
	LastUpdate         time.Time       `yaml:"last_update"`
	Checked            bool            `yaml:"checked"`
	SBOM               SBOM            `yaml:"sbom"`
	OSV                []*OSV          `yaml:"osv"`
	Config             Config          `yaml:"config"`
	HasVulnerabilities bool            `yaml:"has_vulnerabilities"`
	GoModHash          uint64          `yaml:"go_mod_hash"`
}

func NewMetaFile(filename string) *MetaFile {
	f, err := os.Open(filename)
	var meta MetaFile
	if err == nil {
		if err = yaml.NewDecoder(f).Decode(&meta); err == nil {
			if meta.Vulnerabilities != nil {
				meta.Vulnerabilities.Fix()
			}
		}
	}
	meta.filename = filename
	if meta.Vulnerabilities == nil {
		meta.Vulnerabilities = make(Vulnerabilities)
	}
	return &meta
}

// safeGetSBOMModulePath safely retrieves the main module path from SBOM, returning empty string if unavailable
func (m *MetaFile) safeGetSBOMModulePath() string {
	if len(m.SBOM.Modules) > 0 {
		return m.SBOM.Modules[0].Path
	}
	return ""
}

func (m *MetaFile) Save() error {
	m.LastUpdate = time.Now()
	if m.Vulnerabilities == nil {
		m.Vulnerabilities = make(Vulnerabilities)
	}
	m.Vulnerabilities.Fix()
	if modulePath := m.safeGetSBOMModulePath(); modulePath != "" {
		m.PackageName = modulePath
	}
	m.HasVulnerabilities = len(m.Vulnerabilities) > 0
	content, err := yaml.Marshal(m)
	if err == nil {
		err = os.WriteFile(m.filename, content, 0644)
	}
	return err
}
