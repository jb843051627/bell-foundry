package alloy

import "math"

// Profile 描述一种钟型剖面：名义直径到预计质量的经验换算参数。
type Profile struct {
	Code        string  `json:"code"`
	Description string  `json:"description"`
	MassFactor  float64 `json:"mass_factor"` // mass = factor * d^exp（kg，d 单位米）
	DiaExponent float64 `json:"dia_exponent"`
	MinDiaMm    float64 `json:"min_dia_mm"`
	MaxDiaMm    float64 `json:"max_dia_mm"`
}

// 内置剖面表：从教堂钟到手铃的常见规格。
var profiles = []Profile{
	{"C-MAJOR", "church central bell", 1180.0, 2.72, 400, 2200},
	{"S-SWISS", "swiss flemish profile", 1090.0, 2.70, 350, 1800},
	{"T-Bourdon", "bourdon tower bell", 1310.0, 2.78, 500, 2600},
	{"H-HARMONIC", "harmonic concert profile", 1155.0, 2.71, 300, 1600},
	{"G-GOTHIC", "gothic beehive profile", 1240.0, 2.75, 380, 2000},
	{"L-LAP", "lapped carillon bell", 980.0, 2.66, 150, 900},
	{"M-MINI", "small hand bell", 850.0, 2.60, 60, 300},
}

// LookupProfile 按剖面代码查找（大小写不敏感）。
func LookupProfile(code string) (Profile, bool) {
	for _, p := range profiles {
		if equalFold(p.Code, code) {
			return p, true
		}
	}
	return Profile{}, false
}

// Profiles 返回全部内置剖面（拷贝）。
func Profiles() []Profile {
	out := make([]Profile, len(profiles))
	copy(out, profiles)
	return out
}

// EstimateMass 按剖面估算给定直径的钟体质量。
func EstimateMass(p Profile, diaMm float64) (float64, bool) {
	if diaMm < p.MinDiaMm || diaMm > p.MaxDiaMm {
		return 0, false
	}
	d := diaMm / 1000.0
	return math.Round(math.Pow(d, p.DiaExponent)*p.MassFactor*10) / 10, true
}

// EstimateNominalHz 由直径粗估名义频率（Hz），经验公式 2100/d(m) 的幂律修正。
func EstimateNominalHz(diaMm float64) float64 {
	d := diaMm / 1000.0
	if d <= 0 {
		return 0
	}
	hz := 2054.0 * math.Pow(d, -1.02)
	return math.Round(hz*10) / 10
}

func equalFold(a, b string) bool {
	if len(a) != len(b) {
		return false
	}
	for i := 0; i < len(a); i++ {
		ca, cb := a[i], b[i]
		if ca >= 'A' && ca <= 'Z' {
			ca += 'a' - 'A'
		}
		if cb >= 'A' && cb <= 'Z' {
			cb += 'a' - 'A'
		}
		if ca != cb {
			return false
		}
	}
	return true
}
