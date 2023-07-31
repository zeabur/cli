package model

import "time"

type Log struct {
	Timestamp time.Time `json:"timestamp" graphql:"timestamp"`
	Message   string    `json:"message" graphql:"message"`
}

type Logs []*Log

func (l Logs) Header() []string {
	return []string{"Timestamp", "Message"}
}

func (l Logs) Rows() [][]string {
	rows := make([][]string, 0, len(l))
	for _, log := range l {
		rows = append(rows, []string{log.Timestamp.Format(time.RFC3339), log.Message})
	}
	return rows
}

func (l *Log) Header() []string {
	return Logs{l}.Header()
}

func (l *Log) Rows() [][]string {
	return Logs{l}.Rows()
}

var (
	_ Tabler = (*Logs)(nil)
	_ Tabler = (*Log)(nil)
)
