package primitives

type EffectDurationType int

const (
	EffectDurationTypeAllPermanent EffectDurationType = iota
	EffectDurationTypeAllTemporaryHealPermanent
	EffectDurationTypeAllTemporary
)

type EffectDuration struct {
	Type  EffectDurationType
	Steps uint32
}

type Effect struct {
	Duration   *EffectDuration
	Attributes *Attributes
}
