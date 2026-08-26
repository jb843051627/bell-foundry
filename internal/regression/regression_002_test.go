package regression

import ("sync"; "testing"; "github.com/jb843051627/bell-foundry/internal/metrics")

func TestBug02_RegistryConcurrentReadWrite(t *testing.T) {
	registry := metrics.New(); var wg sync.WaitGroup
	for i := 0; i < 8; i++ { wg.Add(1); go func() { defer wg.Done(); for j := 0; j < 1000; j++ { registry.Add("heat", 1); _ = registry.Get("heat") } }() }
	wg.Wait()
	if registry.Get("heat") != 8000 { t.Fatalf("counter lost updates") }
}
