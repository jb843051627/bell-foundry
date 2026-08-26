package regression

import ("sync"; "testing"; "time"; "github.com/jb843051627/bell-foundry/internal/metrics")

func TestBug06_HealthCounterRace(t *testing.T) {
	collector := metrics.NewHealthCollector(); var wg sync.WaitGroup
	for i := 0; i < 8; i++ { wg.Add(1); go func() { defer wg.Done(); for j := 0; j < 200; j++ { collector.ObserveRequest(time.Millisecond, j%2 == 0); _ = collector.Snapshot() } }() }
	wg.Wait()
	if collector.Snapshot()["requests"].(uint64) != 1600 { t.Fatalf("request count lost") }
	collector.Reset()
	if collector.Snapshot()["requests"].(uint64) != 0 { t.Fatalf("reset left stale counter") }
}
