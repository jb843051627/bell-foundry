package regression

import ("testing"; "time"; "github.com/jb843051627/bell-foundry/internal/ingest"; "github.com/jb843051627/bell-foundry/internal/store")

func TestBug26_StreamAndCursorRelease(t *testing.T) { stream := ingest.NewStream(1, nil); closed := make(chan struct{}); go func() { stream.Close(); close(closed) }(); select { case <-closed: case <-time.After(300*time.Millisecond): t.Fatal("sensor stream did not close") }; repo, err := store.Open(t.TempDir()+"/foundry.db"); if err != nil { t.Fatal(err) }; defer repo.Close(); _ = repo.Save("heat", "h1", map[string]string{"status":"ready"}); _ = repo.List("heat", func([]byte) error { return nil }); second := make(chan error, 1); go func() { second <- repo.List("heat", func([]byte) error { return nil }) }(); select { case err := <-second: if err != nil { t.Fatal(err) }; case <-time.After(300*time.Millisecond): t.Fatal("cursor remained open") } }
