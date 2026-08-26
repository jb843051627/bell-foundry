package regression

import "testing"
import "github.com/jb843051627/bell-foundry/internal/service"

func TestBug18_BellDefectReference(t *testing.T) {
	for _, resolve := range []func(any) (string, error){
		func(_ any) (string, error) { return service.ResolveBellReference(nil) },
		func(_ any) (string, error) { return service.ResolveDefectReference(nil) },
	} { if _, err := resolve(nil); err == nil { t.Fatal("missing entity was accepted") } }
}
