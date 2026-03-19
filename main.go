package main

import (
	"Tamelien/blog-aggregator/internal/config"
	"log"
	"os"
)

func main() {
	cfg, err := config.Read()
	if err != nil {
		log.Fatal(err)
	}
	State := state{cfg: &cfg}

	cmds := commands{handlerFunctions: make(map[string]func(*state, command) error)}
	cmds.register("login", handlerLogin)
	args := os.Args
	if len(args) < 2 {
		log.Fatal("Usage: blog-aggregator <command> [args...]")
	}

	cmdArgs := command{
		name: args[1],
		args: args[2:],
	}

	if err := cmds.run(&State, cmdArgs); err != nil {
		log.Fatal(err)
	}
}
