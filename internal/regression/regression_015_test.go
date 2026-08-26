package regression

import ("context"; "errors"; "testing"; "github.com/jb843051627/bell-foundry/internal/model"; "github.com/jb843051627/bell-foundry/internal/service"; "github.com/jb843051627/bell-foundry/internal/store")

func TestBug15_MissingSpecError(t *testing.T) { repo, err := store.Open(t.TempDir()+"/foundry.db"); if err != nil { t.Fatal(err) }; defer repo.Close(); lab := service.NewLab(repo); _, err = lab.GetSpec(context.Background(), "missing-spec"); if !errors.Is(err, model.ErrNotFound) { t.Fatalf("expected ErrNotFound, got %v", err) } }
