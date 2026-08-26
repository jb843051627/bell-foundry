package regression

import "testing"
import "github.com/jb843051627/bell-foundry/internal/model"
import "github.com/jb843051627/bell-foundry/internal/tuning"

func TestBug22_CoolingAndCentsBounds(t *testing.T) { curve := &model.CoolingCurve{SolidusC: 700}; if curve.BelowSolidus() { t.Fatal("empty curve reported solid") }; if values := tuning.PartialCents([]float64{440, 220}, []float64{440}); len(values) != 1 { t.Fatalf("expected one partial, got %d", len(values)) } }
