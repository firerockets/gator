package main

import (
	"database/sql"
	"fmt"
	"os"

	"github.com/firerockets/gator/internal/config"
	"github.com/firerockets/gator/internal/database"

	_ "github.com/lib/pq"
)

func main() {

	cfg, err := config.Read()

	if err != nil {
		os.Exit(0)
		return
	}

	db, err := sql.Open("postgres", cfg.DbURL)

	if err != nil {
		os.Exit(0)
		return
	}

	state := state{
		db:     database.New(db),
		config: &cfg,
	}

	args := os.Args

	if len(args) < 2 {
		fmt.Println("Expected arguments")
		os.Exit(1)
	}

	cmd := command{
		name: args[1],
		args: args[2:],
	}

	cmds := registerCommands()

	err = cmds.run(&state, cmd)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
}

func registerCommands() commands {
	cmds := commands{
		dict: make(map[string]func(*state, command) error),
	}

	cmds.register("login", middlewareLoggedIn(handlerLogin))
	cmds.register("register", middlewareLoggedIn(handlerRegister))
	cmds.register("reset", middlewareLoggedIn(handlerReset))
	cmds.register("users", middlewareLoggedIn(handlerUsers))
	cmds.register("agg", middlewareLoggedIn(handlerAgg))
	cmds.register("addfeed", middlewareLoggedIn(handlerAddFeed))
	cmds.register("feeds", middlewareLoggedIn(handlerFeeds))
	cmds.register("follow", middlewareLoggedIn(handleFollow))
	cmds.register("following", middlewareLoggedIn(handleFollowing))
	cmds.register("unfollow", middlewareLoggedIn(handleUnfollow))
	cmds.register("browse", middlewareLoggedIn(handlerBrowse))

	return cmds
}
