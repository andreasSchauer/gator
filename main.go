package main

import (
	"fmt"
	"log"
	"os"

	"github.com/andreasSchauer/gator/internal/config"
)

type state struct {
	cfg		*config.Config
}


func main() {
	cfg, err := config.Read()
	if err != nil {
		log.Fatalf("error reading config: %v", err)
	}

	appState := &state{
		cfg: &cfg,
	}

	cmds := commands{
		registeredCommands: make(map[string]func(*state, command) error),
	}

	cmds.register("login", handlerLogin)

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