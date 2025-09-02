package check

import (
	"encoding/json"
	"strings"
)

type Config struct {
	ProtocolVersion string `json:"protocol_version"`
	ScannerName     string `json:"scanner_name"`
	ScannerVersion  string `json:"scanner_version"`
	DB              string `json:"db"`
	DBLastModified  string `json:"db_last_modified"`
	GoVersion       string `json:"go_version"`
	ScanLevel       string `json:"scan_level"`
	ScanMode        string `json:"scan_mode"`
}

func GetConfig(data string) *Config {
	if !strings.Contains(data[:20], `"config":`) {
		return nil
	}
	var config struct {
		CONFIG Config `json:"config"`
	}
	if err := json.Unmarshal([]byte(data), &config); err != nil {
		return nil
	}
	return &config.CONFIG
}

/*
{
  "config": {
    "protocol_version": "v1.0.0",
    "scanner_name": "govulncheck",
    "scanner_version": "v1.1.4",
    "db": "https://vuln.go.dev",
    "db_last_modified": "2025-08-14T18:26:33Z",
    "go_version": "go1.24.6",
    "scan_level": "symbol",
    "scan_mode": "source"
  }
}
*/
