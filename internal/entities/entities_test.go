package entities

import (
	_ "embed"
	"testing"

	"github.com/stretchr/testify/assert"
)

//go:embed osv.json
var osv string

//go:embed sbom.json
var sbom string

//go:embed config.json
var config string

func TestGetOSV(t *testing.T) {
	t.Run("osv_generic_empty", func(t *testing.T) {
		osv := get[OSV]("osv", "")
		assert.Nil(t, osv)
	})
	t.Run("empty_osv", func(t *testing.T) {
		osv := GetOSV("")
		assert.Nil(t, osv)
	})
	t.Run("invalid_osv", func(t *testing.T) {
		osv := GetOSV(`"osv":`)
		assert.Nil(t, osv)
	})
	t.Run("valid_osv", func(t *testing.T) {
		osv := GetOSV(osv)
		assert.NotNil(t, osv)
		assert.Equal(t, "GO-2022-0603", osv.Id)

		for _, a := range osv.Affected {
			assert.NotPanics(t, func() {
				_, _ = a.Versions()
			})
			break
		}
	})

}

func TestGetSBOM(t *testing.T) {
	sbom := GetSBOM(sbom)
	assert.NotNil(t, sbom)
	assert.Equal(t, "go1.24.6", sbom.GoVersion)
}

func TestGetConfig(t *testing.T) {
	config := GetConfig(config)
	assert.NotNil(t, config)
	assert.Equal(t, "v1.0.0", config.ProtocolVersion)
}
