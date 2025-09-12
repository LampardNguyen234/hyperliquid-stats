package common

import (
	"fmt"
	"github.com/olekukonko/tablewriter"
	tww "github.com/olekukonko/tablewriter/tw"
)

type writer struct {
	data []byte
}

func newWriter() *writer {
	data := make([]byte, 0)
	return &writer{data: data}
}

func (w *writer) Write(data []byte) (int, error) {
	w.data = append(w.data, data...)
	return len(data), nil
}

type TableFormatter struct {
	*tablewriter.Table
	w *writer
}

func NewTableFormatter() *TableFormatter {
	w := newWriter()
	table := tablewriter.NewWriter(w)

	return &TableFormatter{
		Table: table,
		w:     w,
	}
}

func (tw *TableFormatter) Flush() {
	w := newWriter()
	tw.Table = tablewriter.NewWriter(w)
	tw.w = w
	return
}

func (tw *TableFormatter) WithHeader(fields ...string) *TableFormatter {
	tw.Table.Header(fields)
	return tw
}

func (tw *TableFormatter) WithFooter(fields ...string) *TableFormatter {
	tw.Table.Footer(fields)
	return tw
}

func (tw *TableFormatter) WithCaption(caption string) *TableFormatter {
	tw.Table.Caption(tww.Caption{
		Text: caption,
	})
	return tw
}

func (tw *TableFormatter) WithRow(fields ...interface{}) *TableFormatter {
	row := make([]string, len(fields))
	for i, v := range fields {
		row[i] = fmt.Sprintf("%v", v)
	}
	tw.Table.Append(row)

	return tw
}

func (tw *TableFormatter) WithRows(rows ...[]interface{}) *TableFormatter {
	for _, row := range rows {
		tw.WithRow(row...)
	}

	return tw
}

func (tw *TableFormatter) String() string {
	if len(tw.w.data) == 0 {
		tw.Table.Render()
	}

	return string(tw.w.data)
}
