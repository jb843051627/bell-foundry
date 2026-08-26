package main

import (
	"log"
	"net/http"
	"os"

	"github.com/jb843051627/bell-foundry/internal/handler"
	"github.com/jb843051627/bell-foundry/internal/service"
	"github.com/jb843051627/bell-foundry/internal/store"
)

func main() {
	path := os.Getenv("BELL_FOUNDRY_DB")
	if path == "" {
		path = "data/bell-foundry.db"
	}
	repository, err := store.Open(path)
	if err != nil {
		log.Fatal(err)
	}
	defer repository.Close()
	app := service.NewLab(repository)
	defer app.Close()
	addr := os.Getenv("BELL_FOUNDRY_ADDR")
	if addr == "" {
		addr = ":8080"
	}
	log.Printf("bell-foundry listening on %s", addr)
	log.Fatal(http.ListenAndServe(addr, handler.New(app)))
}
