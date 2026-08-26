package regression

import "testing"
import "github.com/jb843051627/bell-foundry/internal/alloy"
import "github.com/jb843051627/bell-foundry/internal/service"

func TestBug29_RecipeAndTagState(t *testing.T) {
	source := map[string]float64{"copper": 80, "tin": 20}
	plan := alloy.NormalizePlan(source, 100)
	plan["copper"] = 1
	if source["copper"] != 80 { t.Fatal("recipe snapshot aliases input state") }
	labels := []string{"drying"}
	snapshot := service.SnapshotTags(labels)
	snapshot[0] = "closed"
	if labels[0] != "drying" { t.Fatal("tag snapshot aliases state") }
}
