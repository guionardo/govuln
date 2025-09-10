package output

import (
	"bytes"
	"testing"

	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	w := bytes.NewBufferString("")
	o := New(w, "test", true)
	o.AppendHeader(table.Row{"Id", "Name"})
	o.AppendRow(table.Row{1, "Guionardo"})
	o.AppendSeparator()
	o.AppendRow(table.Row{2, "Marines"})
	o.AppendRow(table.Row{3, "João"})
	o.AppendRow(table.Row{4, "Benjamin"})
	o.AppendFooter(table.Row{"Finish"})
	o.Render()

	assert.Greater(t, w.Len(), 0)

}
