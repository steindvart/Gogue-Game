package entity

type PortalType uint

const (
	PortalTypeStart PortalType = iota
	PortalTypeFinish
)

type Portal struct {
	Shape      Box
	PortalType PortalType
}
