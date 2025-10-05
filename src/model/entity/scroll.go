package entity

type Scroll struct {
	Geometry           Object
	AffectedAttributes Attributes[uint]
	Increment          uint
	Name               string
}
