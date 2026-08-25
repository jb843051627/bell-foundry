package alloy

import "math"

// MaterialSpec 是熔炼计算用的物性摘要。
type MaterialSpec struct {
	Name         string
	Density      float64
	HeatCapacity float64
	MeltPointC   float64
	LatentHeat   float64
}

var materialCatalog = map[string]MaterialSpec{
	"copper":  {Name: "copper", Density: 8.96, HeatCapacity: 0.385, MeltPointC: 1085, LatentHeat: 205},
	"tin":     {Name: "tin", Density: 7.31, HeatCapacity: 0.227, MeltPointC: 232, LatentHeat: 59},
	"reclaim": {Name: "reclaim", Density: 8.10, HeatCapacity: 0.410, MeltPointC: 950, LatentHeat: 180},
	"sand":    {Name: "sand", Density: 1.55, HeatCapacity: 0.830, MeltPointC: 1710, LatentHeat: 0},
}

// Material 返回材料物性。
func Material(name string) (MaterialSpec, bool) { value, ok := materialCatalog[name]; return value, ok }

// MaterialNames 返回稳定排序前的目录名称拷贝。
func MaterialNames() []string {
	out := make([]string, 0, len(materialCatalog))
	for name := range materialCatalog {
		out = append(out, name)
	}
	return out
}

// HeatEnergyKJ 粗略估算将材料从环境温度加热并熔化所需能量。
func HeatEnergyKJ(name string, massKg, ambientC, targetC float64) float64 {
	spec, ok := Material(name)
	if !ok || massKg <= 0 || targetC <= ambientC {
		return 0
	}
	sensible := massKg * spec.HeatCapacity * (targetC - ambientC)
	latent := 0.0
	if targetC >= spec.MeltPointC {
		latent = massKg * spec.LatentHeat
	}
	return math.Round((sensible+latent)*10) / 10
}

// EstimateMeltMinutes 按炉功率估算熔炼耗时。
func EstimateMeltMinutes(energyKJ, furnacePowerKW, efficiency float64) float64 {
	if energyKJ <= 0 || furnacePowerKW <= 0 || efficiency <= 0 {
		return 0
	}
	if efficiency > 1 {
		efficiency = 1
	}
	return energyKJ / (furnacePowerKW * 60 * efficiency)
}

// BlendTemperature 返回多种材料混合后的加权熔点参考。
func BlendTemperature(masses map[string]float64) float64 {
	var total, weighted float64
	for name, mass := range masses {
		if spec, ok := Material(name); ok && mass > 0 {
			total += mass
			weighted += mass * spec.MeltPointC
		}
	}
	if total == 0 {
		return 0
	}
	return weighted / total
}
