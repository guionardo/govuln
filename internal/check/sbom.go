package check

import (
	"encoding/json"
	"strings"
)

type SBOM struct {
	GoVersion string   `json:"go_version"`
	Modules   []Module `json:"modules"`
	Roots     []string `json:"roots"`
}

type Module struct {
	Path    string `json:"path"`
	Version string `json:"version,omitempty"`
}

func GetSBOM(data string) *SBOM {
	if !strings.Contains(data[:10], `"SBOM":`) {
		return nil
	}
	var sbom struct {
		SBOM SBOM `json:"SBOM"`
	}
	if err := json.Unmarshal([]byte(data), &sbom); err != nil {
		return nil
	}
	return &sbom.SBOM
}
