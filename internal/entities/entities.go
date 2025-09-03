package entities

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/hashicorp/go-version"
)

type (
	OSV struct {
		Id        string     `json:"id"`
		Modified  string     `json:"modified"`
		Published string     `json:"published"`
		Summary   string     `json:"summary"`
		Affected  []Affected `json:"affected"`
	}
	Affected struct {
		Package Package `json:"package"`
		Ranges  []Range `json:"ranges"`
	}
	Package struct {
		Name      string `json:"name"`
		Ecosystem string `json:"ecosystem"`
	}
	Range struct {
		Type   string  `json:"type"`
		Events []Event `json:"events"`
	}
	Event struct {
		Introduced string `json:"introduced,omitempty"`
		Fixed      string `json:"fixed,omitempty"`
	}

	SBOM struct {
		GoVersion string   `json:"go_version"`
		Modules   []Module `json:"modules"`
		Roots     []string `json:"roots"`
	}
	Module struct {
		Path    string `json:"path"`
		Version string `json:"version,omitempty"`
	}
	Config struct {
		ProtocolVersion string `json:"protocol_version"`
		ScannerName     string `json:"scanner_name"`
		ScannerVersion  string `json:"scanner_version"`
		DB              string `json:"db"`
		DBLastModified  string `json:"db_last_modified"`
		GoVersion       string `json:"go_version"`
		ScanLevel       string `json:"scan_level"`
		ScanMode        string `json:"scan_mode"`
	}
)

func get[T any](prefix, data string) *T {
	prefix = fmt.Sprintf(`"%s":`, prefix)
	if strings.Contains(data[:min(len(data), 15)], prefix) {
		var st struct {
			D T `json:"entity"`
		}

		if err := json.Unmarshal([]byte(strings.Replace(data, prefix, `"entity":`, 1)), &st); err == nil {
			return &st.D
		}
	}
	return nil
}

func GetOSV(data string) *OSV {
	return get[OSV]("osv", data)
}

func (a *Affected) Versions() (introduced *version.Version, fixed *version.Version) {
	for _, r := range a.Ranges {
		for _, e := range r.Events {
			if e.Introduced != "" {
				if i, err := version.NewVersion(e.Introduced); err == nil && (introduced == nil || i.LessThan(introduced)) {
					introduced = i
				}
			}
			if e.Fixed != "" {
				if f, err := version.NewVersion(e.Fixed); err == nil && (fixed == nil || f.GreaterThan(fixed)) {
					fixed = f
				}
			}
		}
	}
	return
}

func GetSBOM(data string) *SBOM {
	return get[SBOM]("SBOM", data)
}

func GetConfig(data string) *Config {
	return get[Config]("config", data)
}
