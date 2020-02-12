package table2

// Cell ...
type Cell struct {
	value    interface{}
	baseType BaseType
}

// NewCell ...
func NewCell(value interface{}) Cell {
	return Cell{value: value, baseType: baseTypeOf(value)}
}
