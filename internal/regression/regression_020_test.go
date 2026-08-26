package regression

import "testing"
import "github.com/jb843051627/bell-foundry/internal/service"

func TestBug20_SensorAndPartialBounds(t *testing.T) { if _, err := service.ReadSampleAt(nil, 0); err == nil { t.Fatal("empty sample accepted") }; if _, err := service.ReadPartialAt([]float64{1}, 2); err == nil { t.Fatal("short partial list accepted") } }
