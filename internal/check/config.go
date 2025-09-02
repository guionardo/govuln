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
