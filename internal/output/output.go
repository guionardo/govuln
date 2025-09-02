package output

import (
	"io"

	"github.com/guionardo/govuln/internal/params"
	"github.com/jedib0t/go-pretty/v6/table"
)

type (
	Output struct {
		w        io.Writer
		t        table.Writer
		markDown bool
	}
	OutputStyle byte
)

const (
	SuccessOutput OutputStyle = iota
	ErrorOutput
)

func New(w io.Writer, title string, withErrorStyle bool) *Output {
	tw := table.NewWriter()
	tw.SetTitle(title)
	o := &Output{
		w:        w,
		t:        tw,
		markDown: params.Parameters.OutputType == "markdown",
	}
	o.SetWithError(withErrorStyle)
	return o
}

func (o *Output) SetWithError(withErrorStyle bool) {
	if params.Parameters.OutputType == "color" {
		if withErrorStyle {
			o.t.SetStyle(table.StyleColoredBlackOnRedWhite)
		} else {
			o.t.SetStyle(table.StyleColoredBlackOnGreenWhite)
		}
	} else {
		if withErrorStyle {
			o.t.SetStyle(table.StyleBold)
		} else {
			o.t.SetStyle(table.StyleRounded)
		}
	}
}
func (o *Output) AppendHeader(row table.Row, configs ...table.RowConfig) {
	o.t.AppendHeader(row, configs...)
}

func (o *Output) AppendSeparator() {
	o.t.AppendSeparator()
}

func (o *Output) AppendRow(row table.Row, configs ...table.RowConfig) {
	o.t.AppendRow(row, configs...)
}

func (o *Output) AppendFooter(row table.Row, configs ...table.RowConfig) {
	o.t.AppendFooter(row, configs...)
}

func (o *Output) Render() {
	var out []byte
	if o.markDown {
		out = []byte(o.t.RenderMarkdown())
	} else {
		out = []byte(o.t.Render())
	}
	_, _ = o.w.Write(append(out, byte('\n')))
}
