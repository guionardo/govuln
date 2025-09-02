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

/*
{
  "SBOM": {
    "go_version": "go1.24.6",
    "modules": [
      {
        "path": "github.com/melisource/govulncheck"
      },
      {
        "path": "github.com/guionardo/go",
        "version": "v0.0.1"
      },
      {
        "path": "gopkg.in/yaml.v3",
        "version": "v3.0.1"
      },
      {
        "path": "stdlib",
        "version": "v1.24.6"
      }
    ],
    "roots": [
      "github.com/melisource/govulncheck",
      "github.com/melisource/govulncheck/internal/exec",
      "github.com/melisource/govulncheck/internal/check",
      "github.com/melisource/govulncheck/internal/store",
      "github.com/melisource/govulncheck/internal/git"
    ]
  }
}*/

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
