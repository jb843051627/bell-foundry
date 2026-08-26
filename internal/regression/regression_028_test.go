package regression

import ("context"; "errors"; "testing"; "github.com/jb843051627/bell-foundry/internal/clock"; "github.com/jb843051627/bell-foundry/internal/notify"; "github.com/jb843051627/bell-foundry/internal/service"; "github.com/jb843051627/bell-foundry/internal/store")

func TestBug28_CanceledNoteAndAlert(t *testing.T) { repo, err := store.Open(t.TempDir()+"/foundry.db"); if err != nil { t.Fatal(err) }; defer repo.Close(); lab := service.NewLabWith(repo, clock.System{}, &notify.MemorySink{}); ctx, cancel := context.WithCancel(context.Background()); cancel(); if _, err := lab.RecordOperatorNote(ctx, "mold", "m1", "operator", "cancelled"); !errors.Is(err, context.Canceled) { t.Fatalf("note accepted after cancel: %v", err) }; sink := &notify.MemorySink{}; if err := sink.Send(ctx, notify.Message{Subject:"cancelled"}); !errors.Is(err, context.Canceled) { t.Fatalf("alert accepted after cancel: %v", err) } }
