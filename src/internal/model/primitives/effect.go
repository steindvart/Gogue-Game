package primitives

type Effect struct {
	// if duration is 0, the effect is instant
	Duration   uint32 // in steps
	Attributes *Attributes
}
