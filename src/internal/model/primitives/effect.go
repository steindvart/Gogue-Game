package primitives

type EffectDurationType int

const (
	EffectDurationTypePermanent EffectDurationType = iota
	EffectDurationTypeTemporary
)

type EffectDuration struct {
	Type  EffectDurationType
	Steps uint32
}

type Effect struct {
	Duration   *EffectDuration
	Attributes *Attributes
}
