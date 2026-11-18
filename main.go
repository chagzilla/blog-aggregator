package main

import (
	"internal/config"
	"internal/database"
	"log"
	"os"
	_ "github.com/lib/pq"
	"database/sql"
	"context"

)


func middlewareLoggedIn(handle func(s *state, cmd command, user database.User) error) func(*state, command) error {
	return func(s *state, c command) error { 
		u, err := s.db.GetUser(context.Background(), s.config.CurrentUserName)
		if err != nil {
			return err
		}
		return handle(s, c, u) 
	}
}

type state struct {
	config * config.Config
	db *database.Queries
}

func main() {
	c := config.Read()
	db, err := sql.Open("postgres", c.DBUrl)
	if err != nil {
		log.Fatal(err)
	}


	s := state { config: &c, db: database.New(db) }

	callbackMap := make(map[string]func(*state, command)error)
	coms := commands { callbacks: callbackMap }
	coms.register("login", handlerLogin)
	coms.register("register", handlerRegister)
	coms.register("reset", handlerReset)
	coms.register("users", handlerUsers)
	coms.register("agg", handlerFeed)
	coms.register("addfeed", middlewareLoggedIn(handlerAddFeed))
	coms.register("feeds", handlerFeeds)
	coms.register("follow", middlewareLoggedIn(handlerFollowFeed))
	coms.register("following", middlewareLoggedIn(handlerFollowing))
	coms.register("unfollow", middlewareLoggedIn(handleUnfollow))

	args := os.Args
	if len(args) < 2 {
		log.Fatal("Not enough arguments")
	}

	err = coms.run(&s, command { name: args[1], arguments: args[2:] })
	if err != nil {
		log.Fatal(err)
	}
}
