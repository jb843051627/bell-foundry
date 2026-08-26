package model

import "time"

// 调音状态。
const (
	TuningUnmeasured  = "unmeasured"
	TuningInTune      = "in_tune"
	TuningNeedsRetune = "needs_retune"
	TuningRecast      = "recast"
)

// Bell 是浇注冷却后的一口钟，记录五分音实测频率。
type Bell struct {
	ID            string    `json:"id"`
	PourID        string    `json:"pour_id"`
	MoldID        string    `json:"mold_id"`
	MassKg        float64   `json:"mass_kg"`
	DiameterMm    float64   `json:"diameter_mm"`
	NominalFreqHz float64   `json:"nominal_freq_hz"`
	Partials      []float64 `json:"partials"` // hum, prime, tierce, quint, nominal（Hz）
	TuningStatus  string    `json:"tuning_status"`
	DetuneCents   float64   `json:"detune_cents"`
	CastAt        time.Time `json:"cast_at"`
}

// 分音槽位下标。
const (
	PartialHum = iota
	PartialPrime
	PartialTierce
	PartialQuint
	PartialNominal
	PartialCount
)

// HasAllPartials 五个分音是否全部测齐。
func (b *Bell) HasAllPartials() bool {
	return len(b.Partials) == PartialCount
}

// NominalMeasured 返回名义分音实测频率；未测齐返回 false。
func (b *Bell) NominalMeasured() (float64, bool) {
	if !b.HasAllPartials() {
		return 0, false
	}
	return b.Partials[PartialNominal], true
}
