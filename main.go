package main

import (
	"database/sql"
	"github.com/acemouty/gator/internal/config"
	"github.com/acemouty/gator/internal/database"
	_ "github.com/lib/pq"
	"log"
	"os"
)

type command struct {
	name string
	args []string
}

type state struct {
	db  *database.Queries
	cfg *config.Config
}

func main() {
	// excludes program name
	args := os.Args[1:]
	if len(args) < 1 {
		log.Fatalf("Expected atleast 1 argument got %v", len(args))
	}

	cfg := config.Read()
	db := dbOpen(cfg.DbUrl)
	dbQueries := database.New(db)

	appState := state{cfg: &cfg, db: dbQueries}
	command := command{name: args[0], args: args[1:]}
	commandStore := commandStore{commandsMap: make(commandMap)}

	commandStore.register("login", handlerLogin)
	commandStore.register("register", handlerRegitser)

	err := commandStore.run(&appState, command)
	if err != nil {
		log.Fatalf("main: encountered and error running '%v': %v", command.name, err)
	}

}

func dbOpen(dbUrl string) *sql.DB {
	db, err := sql.Open("postgres", dbUrl)
	if err != nil {
		log.Fatalf("main: Issue connecting to database: %v", err)
	}

	return db
}
