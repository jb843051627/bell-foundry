package regression

import ("context"; "errors"; "testing"; "github.com/jb843051627/bell-foundry/internal/service")

func TestBug09_CoolingPourCancellation(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background()); cancel(); called := false
	for _, run := range []func(context.Context, func(context.Context) error) error{ service.CoolingReadContext, service.PourReadContext } {
		err := run(ctx, func(context.Context) error { called = true; return nil })
		if !errors.Is(err, context.Canceled) { t.Fatalf("canceled context returned %v", err) }
	}
	if called { t.Fatal("canceled callback ran") }
}
