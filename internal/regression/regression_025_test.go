package regression

import ("context"; "testing"; "time"; "github.com/jb843051627/bell-foundry/internal/notify"; "github.com/jb843051627/bell-foundry/internal/store")

func TestBug25_AlertListResourceRelease(t *testing.T) { sink := &notify.MemorySink{}; _ = sink.Messages(); done := make(chan struct{}); go func() { _ = sink.Messages(); close(done) }(); select { case <-done: case <-time.After(300*time.Millisecond): t.Fatal("alert mutex was not released") }; repo, err := store.Open(t.TempDir()+"/foundry.db"); if err != nil { t.Fatal(err) }; defer repo.Close(); _ = repo.Save("heat", "h1", map[string]string{"status":"ready"}); _ = repo.List("heat", func([]byte) error { return nil }); second := make(chan error, 1); go func() { second <- repo.List("heat", func([]byte) error { return nil }) }(); select { case err := <-second: if err != nil { t.Fatal(err) }; case <-time.After(300*time.Millisecond): t.Fatal("rows cursor was not released") }; _ = context.Background() }
