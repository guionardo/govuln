package check

import (
	"fmt"
	"io"
	"os"
	"slices"
	"strings"
	"sync"

	"github.com/guionardo/govuln/internal/output"
	goversion "github.com/hashicorp/go-version"
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/jedib0t/go-pretty/v6/text"
)

type (
	Vulnerabilities map[string][]*Vulnerability // map[packageName][]*Vulnerability
	Vulnerability   struct {
		Id         string             `yaml:"id"`
		Summary    string             `yaml:"summary"`
		SIntro     string             `yaml:"introduced"`
		SFixed     string             `yaml:"fixed"`
		SCurrent   string             `yaml:"current"`
		Introduced *goversion.Version `yaml:"-"`
		Fixed      *goversion.Version `yaml:"-"`
		Current    *goversion.Version `yaml:"-"`
	}
	SubmodulesVulnerabilities struct {
		sub  map[string]map[string][]string
		lock sync.RWMutex
	}
)

func (v Vulnerabilities) Table(title string, version string) {
	isInErrorStyle := len(v) > 0
	o := output.New(os.Stdout, title, isInErrorStyle)
	if isInErrorStyle {
		o.AppendHeader(table.Row{"Package", "Vulnerability", "Introduced", "Fixed", "Current"})
	}

	previousPackageLength := 0
	for pack, vulns := range v {
		if previousPackageLength > 2 {
			o.AppendSeparator()
		}
		for _, vuln := range vulns {
			o.AppendRow(table.Row{pack, vuln.Id, vuln.Introduced, vuln.Fixed, vuln.Current},
				table.RowConfig{AutoMerge: true, AutoMergeAlign: text.AlignAuto})
		}
		previousPackageLength = len(vulns)
	}

	if len(v) == 0 {
		o.AppendFooter(table.Row{"SUCCESS!! No Vulnerabilities Found", version})
	} else {
		o.AppendFooter(table.Row{"Total", len(v), "", "", version})
	}
	o.Render()
}

func (v *Vulnerability) Fix() {
	if v.Introduced == nil && len(v.SIntro) > 0 {
		if f, err := goversion.NewVersion(v.SIntro); err == nil {
			v.Introduced = f
		}
	} else if v.Introduced != nil && len(v.SIntro) == 0 {
		v.SIntro = v.Introduced.Original()
	}
	if v.Fixed == nil && len(v.SFixed) > 0 {
		if f, err := goversion.NewVersion(v.SFixed); err == nil {
			v.Fixed = f
		}
	} else if v.Fixed != nil && len(v.SFixed) == 0 {
		v.SFixed = v.Fixed.Original()
	}
	if v.Current == nil && len(v.SCurrent) > 0 {
		if f, err := goversion.NewVersion(v.SCurrent); err == nil {
			v.Current = f
		}
	} else if v.Current != nil && len(v.SCurrent) == 0 {
		v.SCurrent = v.Current.Original()
	}
}

func (v *Vulnerabilities) Fix() {
	for _, vulns := range *v {
		for _, vuln := range vulns {
			vuln.Fix()
		}
	}
}

func NewSubmodulesVulnerabilities() *SubmodulesVulnerabilities {
	return &SubmodulesVulnerabilities{
		sub: make(map[string]map[string][]string),
	}
}

func (s *SubmodulesVulnerabilities) Add(module string, version string, vuln string) {
	s.lock.Lock()
	defer s.lock.Unlock()
	if s.sub[module] == nil {
		s.sub[module] = make(map[string][]string)
	}
	s.sub[module][version] = append(s.sub[module][version], vuln)
}

func (s *SubmodulesVulnerabilities) Table(w io.Writer, internalOwner string) {
	s.lock.RLock()
	defer s.lock.RUnlock()
	t := output.New(w, "Submodules vulnerabilities", false)
	if len(s.sub) == 0 {
		t.AppendHeader(table.Row{"No submodules from owner", internalOwner})
		t.Render()
		return
	}
	t.AppendHeader(table.Row{"Module", "Version", "Vulnerability"})

	hasVulnerabilities := false
	for module, versions := range s.sub {
		// w := strings.Split(module, path.Join(internalOwner, ""))
		// module = w[len(w)-1]
		for version, vulns := range versions {
			if !slices.Contains(vulns, "SAFE") {
				hasVulnerabilities = true

			}
			if len(vulns) > 2 {
				vulns = append(vulns[0:2], fmt.Sprintf("+%d", len(vulns)-2))
			}
			t.AppendRow(table.Row{module, version, strings.Join(vulns, " ")})
		}
	}
	if hasVulnerabilities {
		t.SetWithError(true)
		t.AppendFooter(table.Row{"Total", len(s.sub)}, table.RowConfig{AutoMerge: true, AutoMergeAlign: text.AlignCenter})
	} else {
		t.AppendFooter(table.Row{"SUCCESS!! No Vulnerabilities Found"}, table.RowConfig{AutoMerge: true, AutoMergeAlign: text.AlignCenter})
	}
	t.Render()
}
