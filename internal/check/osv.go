package check

import (
	"encoding/json"
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
)

func GetOSV(data string) *OSV {
	if !strings.Contains(data[:10], `"osv":`) {
		return nil
	}
	var osv struct {
		OSV OSV `json:"osv"`
	}
	if err := json.Unmarshal([]byte(data), &osv); err != nil {
		return nil
	}
	return &osv.OSV
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

/*
{
  "osv": {
    "schema_version": "1.3.1",
    "id": "GO-2025-3563",
    "modified": "2025-04-08T19:46:23Z",
    "published": "2025-04-08T19:46:23Z",
    "aliases": [
      "CVE-2025-22871"
    ],
    "summary": "Request smuggling due to acceptance of invalid chunked data in net/http",
    "details": "The net/http package improperly accepts a bare LF as a line terminator in chunked data chunk-size lines. This can permit request smuggling if a net/http server is used in conjunction with a server that incorrectly accepts a bare LF as part of a chunk-ext.",
    "affected": [
      {
        "package": {
          "name": "stdlib",
          "ecosystem": "Go"
        },
        "ranges": [
          {
            "type": "SEMVER",
            "events": [
              {
                "introduced": "0"
              },
              {
                "fixed": "1.23.8"
              },
              {
                "introduced": "1.24.0-0"
              },
              {
                "fixed": "1.24.2"
              }
            ]
          }
        ],
        "ecosystem_specific": {
          "imports": [
            {
              "path": "net/http/internal",
              "symbols": [
                "chunkedReader.Read",
                "readChunkLine"
              ]
            }
          ]
        }
      }
    ],
    "references": [
      {
        "type": "FIX",
        "url": "https://go.dev/cl/652998"
      },
      {
        "type": "REPORT",
        "url": "https://go.dev/issue/71988"
      },
      {
        "type": "WEB",
        "url": "https://groups.google.com/g/golang-announce/c/Y2uBTVKjBQk"
      }
    ],
    "credits": [
      {
        "name": "Jeppe Bonde Weikop"
      }
    ],
    "database_specific": {
      "url": "https://pkg.go.dev/vuln/GO-2025-3563",
      "review_status": "REVIEWED"
    }
  }
}*/
