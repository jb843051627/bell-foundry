package regression

import "testing"
import "github.com/jb843051627/bell-foundry/internal/notify"
import "github.com/jb843051627/bell-foundry/internal/service"
import "github.com/jb843051627/bell-foundry/internal/store"
import "context"

func TestBug24_PartialsAndAlertsSnapshot(t *testing.T) { source := []string{"hum"}; copy := service.SnapshotPartials(source); copy[0] = "changed"; if source[0] != "hum" { t.Fatal("partials snapshot mutated source") }; sink := &notify.MemorySink{}; if err := sink.Send(context.Background(), notify.Message{Body: "original"}); err != nil { t.Fatal(err) }; messages := sink.Messages(); messages[0].Body = "changed"; if sink.Messages()[0].Body != "original" { t.Fatal("alert snapshot mutated sink") }; _ = store.Store{} }
