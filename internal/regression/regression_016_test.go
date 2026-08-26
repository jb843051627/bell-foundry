package regression

import ("context"; "errors"; "testing"; "github.com/jb843051627/bell-foundry/internal/clock"; "github.com/jb843051627/bell-foundry/internal/model"; "github.com/jb843051627/bell-foundry/internal/notify"; "github.com/jb843051627/bell-foundry/internal/service"; "github.com/jb843051627/bell-foundry/internal/store")

type failingSink struct{ err error }
func (s failingSink) Send(context.Context, notify.Message) error { return s.err }
func TestBug16_AlertSinkError(t *testing.T) { repo, err := store.Open(t.TempDir()+"/foundry.db"); if err != nil { t.Fatal(err) }; defer repo.Close(); sentinel := errors.New("webhook unavailable"); lab := service.NewLabWith(repo, clock.System{}, failingSink{err: sentinel}); _, err = lab.RaiseAlert(context.Background(), "pour:1", "critical", "hot"); if !errors.Is(err, sentinel) { t.Fatalf("sink error lost: %v", err) }; if !errors.Is(service.PersistAlertResult(sentinel), sentinel) { t.Fatal("alert error identity lost") }; _ = model.AlertCritical }
