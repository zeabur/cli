package printer

import (
	"os"

	"github.com/olekukonko/tablewriter"
)

type defaultPrinter struct{}

func New() Printer {
	return &defaultPrinter{}
}

func (p *defaultPrinter) Table(header []string, rows [][]string) {
	columnsCount := len(header)

	colors := []tablewriter.Colors{
		{tablewriter.FgHiMagentaColor},
		{tablewriter.FgGreenColor},
		{tablewriter.FgHiBlueColor},
		{tablewriter.FgHiYellowColor},
		{tablewriter.FgHiCyanColor},
	}
	headerColor := tablewriter.Colors{tablewriter.Bold, tablewriter.FgBlackColor}

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader(header)
	table.SetBorder(false)
	table.SetColumnSeparator("")

	headerColors := make([]tablewriter.Colors, columnsCount)
	for i := 0; i < columnsCount; i++ {
		headerColors[i] = headerColor
	}
	table.SetHeaderColor(headerColors...)

	columnColors := make([]tablewriter.Colors, columnsCount)
	for i := 0; i < columnsCount; i++ {
		columnColors[i] = colors[i%len(colors)]
	}

	table.SetColumnColor(columnColors...)

	// fix take space as newline bug
	table.SetAutoWrapText(false)

	table.AppendBulk(rows)

	table.Render()
}

var _ Printer = &defaultPrinter{}
