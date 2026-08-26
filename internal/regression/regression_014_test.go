package regression

import ("context"; "database/sql"; "errors"; "testing"; "github.com/jb843051627/bell-foundry/internal/model"; "github.com/jb843051627/bell-foundry/internal/store"; "github.com/jb843051627/bell-foundry/internal/service")

func TestBug14_TransactionErrorIdentity(t *testing.T) {
	sentinel := errors.New("transaction-sentinel"); repo, err := store.Open(t.TempDir()+"/foundry.db"); if err != nil { t.Fatal(err) }; defer repo.Close()
	err = repo.Transaction(func(tx *sql.Tx) error { return sentinel }); if !errors.Is(err, sentinel) { t.Fatalf("transaction identity lost: %v", err) }
	if !errors.Is(service.PersistInspectionResult(sentinel), sentinel) { t.Fatal("inspection identity lost") }; _ = context.Background(); _ = model.ErrNotFound
}
