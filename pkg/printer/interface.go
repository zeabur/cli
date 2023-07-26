package printer

type Printer interface {
	Table(header []string, rows [][]string) // Print a table with the given header and rows
}
