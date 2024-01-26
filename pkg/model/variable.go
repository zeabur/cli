package model

type Variable struct {
	Key       string `json:"key" graphql:"key"`
	Value     string `json:"value" graphql:"value"`
	ServiceID string `json:"serviceID" graphql:"serviceID"`
}

type Variables []*Variable

func (v Variables) Header() []string {
	return []string{"Key", "Value"}
}

func (v Variables) Rows() [][]string {
	rows := make([][]string, 0, len(v))
	headerLen := len(v.Header())
	for _, variable := range v {
		row := make([]string, 0, headerLen)
		row = append(row, variable.Key)
		row = append(row, variable.Value)

		rows = append(rows, row)
	}
	return rows
}

func (v Variables) ToMap() map[string]string {
	m := make(map[string]string)
	for _, variable := range v {
		m[variable.Key] = variable.Value
	}
	return m
}
