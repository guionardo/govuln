package store

import (
	"testing"

	goversion "github.com/hashicorp/go-version"
	"github.com/stretchr/testify/assert"
)

func TestVulnerability_Fix(t *testing.T) {
	ver, _ := goversion.NewVersion("1.0.0")

	t.Run("fix_from_versions", func(t *testing.T) {
		v := &Vulnerability{
			Introduced: ver,
			Fixed:      ver,
			Current:    ver,
		}
		v.Fix()
		assert.Equal(t, "1.0.0", v.SIntro)
		assert.Equal(t, "1.0.0", v.SFixed)
		assert.Equal(t, "1.0.0", v.SCurrent)
	})
	t.Run("fix_from_strings", func(t *testing.T) {
		v := &Vulnerability{
			SIntro:   "1.0.0",
			SFixed:   "1.0.0",
			SCurrent: "1.0.0",
		}
		v.Fix()
		assert.Equal(t, ver, v.Introduced)
		assert.Equal(t, ver, v.Fixed)
		assert.Equal(t, ver, v.Current)
	})
}
