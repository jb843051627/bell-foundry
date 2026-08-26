package regression

import ("context"; "testing"; "time"; "github.com/jb843051627/bell-foundry/internal/ingest")

func TestBug27_QueueWaitForWorker(t *testing.T) { queue := ingest.New(1); started := make(chan struct{}); release := make(chan struct{}); if err := queue.Submit(context.Background(), ingest.Job{Run: func(context.Context) error { close(started); <-release; return nil }}); err != nil { t.Fatal(err) }; <-started; closed := make(chan struct{}); go func() { queue.Close(); close(closed) }(); select { case <-closed: t.Fatal("queue close returned before worker finished"); case <-time.After(50*time.Millisecond): }; close(release); select { case <-closed: case <-time.After(500*time.Millisecond): t.Fatal("queue close did not wait") } }
