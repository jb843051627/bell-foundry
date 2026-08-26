package regression

import ("errors"; "testing"; "github.com/jb843051627/bell-foundry/internal/service")

func TestBug11_HeatMoldErrorIdentity(t *testing.T) {
	sentinel := errors.New("HeatMoldErrorIdentity-sentinel")
	for _, wrap := range []func(error) error{ service.PersistHeatResult, service.PersistMoldResult } { if !errors.Is(wrap(sentinel), sentinel) { t.Fatal("persistence error identity lost") } }
}
