package regression

import ("sync"; "testing"; "time"; "github.com/jb843051627/bell-foundry/internal/notify")

func TestBug03_AlertDeduplication(t *testing.T) {
	dedup := notify.NewDeduper(time.Hour); var wg sync.WaitGroup; var mu sync.Mutex; allowed := 0
	for i := 0; i < 32; i++ { wg.Add(1); go func() { defer wg.Done(); if dedup.Allow("pour:critical", time.Unix(1, 0)) { mu.Lock(); allowed++; mu.Unlock() }; _ = dedup.Size() }() }
	wg.Wait(); if allowed != 1 { t.Fatalf("expected one alert, got %d", allowed) }
}
