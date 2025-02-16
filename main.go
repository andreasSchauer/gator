package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/andreasSchauer/gator/internal/config"
	"github.com/andreasSchauer/gator/internal/database"
	_ "github.com/lib/pq"
)

type state struct {
	db  *database.Queries
	cfg *config.Config
}

func main() {
	cfg, err := config.Read()
	if err != nil {
		log.Fatalf("error reading config: %v", err)
	}

	db, err := sql.Open("postgres", cfg.DbURL)
	if err != nil {
		log.Fatalf("error loading database: %v", err)
	}

	appState := &state{
		db:  database.New(db),
		cfg: &cfg,
	}

	cmds := commands{
		registeredCommands: make(map[string]func(*state, command) error),
	}

	cmds.register("login", handlerLogin)
	cmds.register("register", handlerRegister)
	cmds.register("reset", handlerReset)
	cmds.register("users", handlerUsers)
	cmds.register("agg", handlerAgg)
	cmds.register("addfeed", handlerAddFeed)
	cmds.register("feeds", handlerFeeds)
	cmds.register("follow", handlerFollow)
	cmds.register("following", handlerFollowing)

	if len(os.Args) < 2 {
		fmt.Println("Usage: cli <command> [args...]")
		os.Exit(1)
	}

	userCommand := command{
		name: os.Args[1],
		args: os.Args[2:],
	}

	err = cmds.run(appState, userCommand)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
