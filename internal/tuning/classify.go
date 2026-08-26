package tuning

// Result 是一次调音评估的完整结果。
type Result struct {
	Cents          []float64 `json:"cents"`
	Worst          float64   `json:"worst"`
	Mean           float64   `json:"mean"`
	Status         string    `json:"status"`
	WorstName      string    `json:"worst_name"`
	Recommendation string    `json:"recommendation"`
}

// Classify 根据最大音分偏差分级。
func Classify(worst, tuneLimit, retuneLimit float64) string {
	if worst <= tuneLimit {
		return "in_tune"
	}
	if worst <= retuneLimit {
		return "needs_retune"
	}
	return "recast"
}

// Evaluate 对五分音目标与实测值做一次完整评估。
func Evaluate(targets, measured []float64, tuneLimit, retuneLimit float64) Result {
	if !Complete(measured) || len(targets) != len(measured) {
		return Result{Status: "unmeasured", Recommendation: "complete all partial measurements"}
	}
	cents := PartialCents(targets, measured)
	worst := 0.0
	worstIndex := 0
	for i, value := range cents {
		if AbsCents(value) > worst {
			worst = AbsCents(value)
			worstIndex = i
		}
	}
	status := Classify(worst, tuneLimit, retuneLimit)
	recommendation := "release"
	if status == "needs_retune" {
		recommendation = "remove material near the rim and remeasure"
	} else if status == "recast" {
		recommendation = "hold bell for metallurgical review"
	}
	return Result{
		Cents: cents, Worst: worst, Mean: MeanCents(cents), Status: status,
		WorstName: PartialName[worstIndex], Recommendation: recommendation,
	}
}

// RetuneSteps 将偏差换算为简单的工艺调整步骤。
func RetuneSteps(result Result) []string {
	if result.Status != "needs_retune" {
		return nil
	}
	direction := "remove"
	if result.CentsForWorst() < 0 {
		direction = "add mass"
	}
	return []string{"isolate bell", direction + " around " + result.WorstName, "remeasure all partials"}
}

// CentsForWorst 返回最差分音的有符号偏差。
func (r Result) CentsForWorst() float64 {
	for i, name := range PartialName {
		if name == r.WorstName && i < len(r.Cents) {
			return r.Cents[i]
		}
	}
	return 0
}
