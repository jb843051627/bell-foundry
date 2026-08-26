package regression

import "testing"
import "github.com/jb843051627/bell-foundry/internal/model"
import "github.com/jb843051627/bell-foundry/internal/service"
import "github.com/jb843051627/bell-foundry/internal/tuning"

func TestBug21_DefectAndRatioBounds(t *testing.T) { if _, err := service.ReadDefectAt(nil, 0); err == nil { t.Fatal("empty defect list accepted") }; if report := tuning.RatioReport([]float64{440, 220}); report != nil { t.Fatal("short ratio report should be empty") }; _ = model.DefectOpen }
