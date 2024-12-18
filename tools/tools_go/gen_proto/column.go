package gen_proto

/*
   @Author: orbit-w
   @File: row
   @2024 12月 周日 23:02
*/

type Column struct {
	Name       string `json:"name"`
	ChName     string `json:"ch_name"`
	Desc       string `json:"desc"`
	Type       string `json:"type"`
	Permission string `json:"permission"`
	Key        bool   `json:"key"`
}

func newColumn(name, chName, desc, t, permission string, key bool) *Column {
	return &Column{
		Name:       name,
		ChName:     chName,
		Desc:       desc,
		Type:       t,
		Permission: permission,
		Key:        key,
	}
}

type Columns []*Column

func (c Columns) Len() int {
	return len(c)
}

func (c Columns) Less(i, j int) bool {
	return c[i].Name < c[j].Name
}

func (c Columns) GetColumn(name string) *Column {
	for _, column := range c {
		if column.Name == name {
			return column
		}
	}
	return nil
}
