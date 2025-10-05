package entity

type Scroll struct {
	Shape              Box
	AffectedAttributes Attributes
	Increment          uint
	Name               string
}
