package regression

 import ("context"; "testing"; "time"; "github.com/jb843051627/bell-foundry/internal/ingest")

func TestBug07_QueueStopsWorker(t *testing.T) {
	queue := ingest.New(1); done := make(chan error, 1); if err := queue.Submit(context.Background(), ingest.Job{Run: func(_ context.Context) error { return nil }, Done: done}); err != nil { t.Fatal(err) }
	closed := make(chan struct{}); go func() { queue.Close(); close(closed) }(); select { case <-closed: case <-time.After(500 * time.Millisecond): t.Fatal("queue close did not stop worker") }
}
