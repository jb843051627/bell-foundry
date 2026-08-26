package regression

import ("errors"; "testing"; "github.com/jb843051627/bell-foundry/internal/service")

func TestBug13_AlertInspectionErrorIdentity(t *testing.T) {
	sentinel := errors.New("AlertInspectionErrorIdentity-sentinel")
	for _, wrap := range []func(error) error{ service.PersistAlertResult, service.PersistInspectionResult } { if !errors.Is(wrap(sentinel), sentinel) { t.Fatal("persistence error identity lost") } }
}
