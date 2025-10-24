package config

type ElementType int

const (
	TypeNode ElementType = iota + 1
	TypeComponent
	TypeRelation
	TypeLayer
	TypeCluster
	TypeLabel
	TypeView
)

var elementTypes = [...]string{"", "node", "component", "relation", "layer", "cluster", "label", "view"}

func (v ElementType) String() string {
	return elementTypes[v]
}

// Element is ndiag element.
type Element interface {
	ElementType() ElementType
	Id() string
	FullName() string
	DescFilename() string
}
