package check

import (
	"fmt"
	"io"
	"strings"
	"sync"

	goversion "github.com/hashicorp/go-version"
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/jedib0t/go-pretty/v6/text"
	"github.com/melisource/fury_fbm-fiscal-govulncheck/internal/config"
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

func (v Vulnerabilities) Table(w io.Writer, title string, version string) {
	t := table.NewWriter()
	if len(v) == 0 {
		t.SetStyle(table.StyleColoredBlackOnGreenWhite)
	} else {
		t.SetStyle(table.StyleColoredRedWhiteOnBlack)
		t.AppendHeader(table.Row{"Package", "Vulnerability", "Introduced", "Fixed", "Current"})
	}
	t.SetOutputMirror(w)
	previousPackageLength := 0
	for pack, vulns := range v {
		if previousPackageLength > 2 {
			t.AppendSeparator()
		}
		for _, vuln := range vulns {
			t.AppendRow(table.Row{pack, vuln.Id, vuln.Introduced, vuln.Fixed, vuln.Current})
		}
		previousPackageLength = len(vulns)
	}
	if len(title) > 0 {
		t.SetTitle(title)
	}

	if len(v) == 0 {
		t.AppendFooter(table.Row{"SUCCESS!! No Vulnerabilities Found", version}, table.RowConfig{AutoMerge: true, AutoMergeAlign: text.AlignCenter})
	} else {
		t.AppendFooter(table.Row{"Total", len(v), "", "", version}, table.RowConfig{AutoMerge: true, AutoMergeAlign: text.AlignCenter})
	}
	config.Render(t)
	_, _ = w.Write([]byte("\n"))
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

func (s *SubmodulesVulnerabilities) Table(w io.Writer) {
	s.lock.RLock()
	defer s.lock.RUnlock()
	t := table.NewWriter()
	t.SetTitle("Submodules Vulnerabilities")
	t.SetOutputMirror(w)
	if len(s.sub) == 0 {
		t.SetStyle(table.StyleColoredBlackOnGreenWhite)
	} else {
		t.SetStyle(table.StyleColoredRedWhiteOnBlack)
		t.AppendHeader(table.Row{"Module", "Version", "Vulnerability"})
	}

	t.SetStyle(table.StyleColoredRedWhiteOnBlack)
	for module, versions := range s.sub {
		w := strings.Split(module, "melisource/")
		module = w[len(w)-1]
		for version, vulns := range versions {
			if len(vulns) > 2 {
				vulns = append(vulns[0:2], fmt.Sprintf("+%d", len(vulns)-2))
			}
			t.AppendRow(table.Row{module, version, strings.Join(vulns, " ")})
		}
	}
	if len(s.sub) == 0 {
		t.AppendFooter(table.Row{"SUCCESS!! No Vulnerabilities Found"}, table.RowConfig{AutoMerge: true, AutoMergeAlign: text.AlignCenter})
	} else {
		t.AppendFooter(table.Row{"Total", len(s.sub)}, table.RowConfig{AutoMerge: true, AutoMergeAlign: text.AlignCenter})
	}
	config.Render(t)
}
