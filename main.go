package main

import (
	"database/sql"
	"fmt"
	"os"

	"github.com/BlessE1/gator/internal/config"
	"github.com/BlessE1/gator/internal/database"

	_ "github.com/lib/pq"
)

type state struct {
	db     *database.Queries
	config *config.Config
}

func main() {
	// Configuration
	cfg, err := config.Read()
	if err != nil {
		fmt.Println(err)
	}

	db, err := sql.Open("postgres", cfg.DbUrl)
	if err != nil {
		fmt.Println(err)
	}
	dbQueries := database.New(db)

	currState := state{
		db:     dbQueries,
		config: &cfg,
	}

	cmds := commands{
		cmdMap: make(map[string]func(*state, command) error),
	}

	cmds.register("login", handlerLogin)
	cmds.register("register", handlerRegister)
	cmds.register("reset", handlerReset)
	cmds.register("users", handlerGetUsers)
	cmds.register("agg", handlerAgg)
	cmds.register("feeds", handlerGetFeeds)
	cmds.register("addfeed", middlewareLoggedIn(handlerAddFeed))
	cmds.register("follow", middlewareLoggedIn(handlerFollow))
	cmds.register("following", middlewareLoggedIn(handlerFollowing))
	cmds.register("unfollow", middlewareLoggedIn(handlerUnfollow))
	cmds.register("browse", middlewareLoggedIn(handlerBrowse))

	if len(os.Args) < 2 {
		fmt.Println("Not enough arguments were provided")
		fmt.Println(os.Args)
		os.Exit(1)
	}

	cmd := command{
		Name: os.Args[1],
		Args: os.Args[2:len(os.Args)],
	}

	err = cmds.run(&currState, cmd)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
