package util

import (
	"github.com/olekukonko/tablewriter"
	"strings"
)

func BasicTable(header []string, data [][]string, footer []string) string {
	tableString := &strings.Builder{}
	tb := tablewriter.NewWriter(tableString)
	tb.SetAlignment(tablewriter.ALIGN_LEFT)
	tb.SetHeader(header)
	tb.AppendBulk(data)
	tb.SetAutoFormatHeaders(false)
	tb.SetRowLine(true)
	tb.SetFooter(footer)
	tb.Render()
	return tableString.String()
}

func SimpleTable(data [][]string) string {
	tableString := &strings.Builder{}
	tb := tablewriter.NewWriter(tableString)
	tb.SetAlignment(tablewriter.ALIGN_LEFT)
	tb.AppendBulk(data)
	tb.SetAutoMergeCells(true)
	tb.SetRowLine(true)
	tb.Render()
	return tableString.String()
}

func VerticalTable(header []string, data [][]string) string {
	tableString := &strings.Builder{}
	tb := tablewriter.NewWriter(tableString)
	tb.SetAlignment(tablewriter.ALIGN_LEFT)
	tb.SetRowLine(true)

	fieldNum := len(data) + 1
	rowNum := len(header)
	d := make([][]string, rowNum)
	for r := 0; r < rowNum; r++ {
		d[r] = make([]string, fieldNum)
		d[r][0] = header[r]
	}
	for field := 1; field < fieldNum; field++ {
		for row := 0; row < rowNum; row++ {
			d[row][field] = data[field-1][row]
		}
	}
	tb.AppendBulk(d)
	tb.Render()
	return tableString.String()
}
