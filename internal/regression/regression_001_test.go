package regression

import ("sync"; "testing"; "github.com/jb843051627/bell-foundry/internal/model")

func TestBug01_LedgerConcurrentUpdates(t *testing.T) {
	ledger := model.NewLedger()
	var wg sync.WaitGroup
	for i := 0; i < 8; i++ { wg.Add(1); go func() { defer wg.Done(); for j := 0; j < 1000; j++ { ledger.Add("pour", 1); _ = ledger.Value("pour") } }() }
	wg.Wait()
	if got := ledger.Value("pour"); got != 8000 { t.Fatalf("expected 8000 pour events, got %d", got) }
}
