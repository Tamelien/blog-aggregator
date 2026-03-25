package main

import (
	"Tamelien/blog-aggregator/internal/config"
	"Tamelien/blog-aggregator/internal/database"
	"database/sql"
	"log"
	"os"

	_ "github.com/lib/pq"
)

func main() {
	cfg, err := config.Read()
	if err != nil {
		log.Fatal(err)
	}
	State := state{cfg: &cfg}

	db, err := sql.Open("postgres", State.cfg.DBURL)
	if err != nil {
		log.Fatalf("failed to open database %s: %v", State.cfg.DBURL, err)
	}
	dbQueries := database.New(db)
	State.db = dbQueries

	cmds := commands{handlerFunctions: make(map[string]func(*state, command) error)}
	cmds.register("login", handlerLogin)
	cmds.register("register", handlerRegister)
	cmds.register("reset", handlerReset)
	cmds.register("users", handlerGetUsers)
	cmds.register("agg", handlerAgg)
	cmds.register("addfeed", middlewareLoggedIn(handlerAddFeed))
	cmds.register("feeds", handlerGetFeeds)
	cmds.register("follow", middlewareLoggedIn(handlerFollow))
	cmds.register("following", middlewareLoggedIn(handlerFollowing))
	cmds.register("unfollow", middlewareLoggedIn(handlerUnfollow))
	cmds.register("browse", middlewareLoggedIn(handlerBrowse))

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
