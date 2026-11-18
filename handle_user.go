package main

import (
	"errors"
	"fmt"
	"internal/database"
	"github.com/google/uuid"
	"time"
	"context"
)


func handlerLogin(s *state, cmd command) error {
	if len(cmd.arguments) == 0 {
		return errors.New("Not enough arguements: Login requires positional username")
	}

	_, err := s.db.GetUser(context.Background(), cmd.arguments[0])
	if err != nil {
		return err
	}
	s.config.SetUser(cmd.arguments[0])
	fmt.Printf("User has been set\n")

	return nil
}

func handlerRegister(s *state, cmd command) error {
	if len(cmd.arguments) == 0 {
		return errors.New("Not enough arguements: Register requires positional username")
	}

	user, err := s.db.CreateUser(context.Background(), database.CreateUserParams { ID: uuid.New(), Name: cmd.arguments[0], CreatedAt: time.Now(), UpdatedAt: time.Now() })
	if err != nil {
		return err
	}
	s.config.SetUser(user.Name)
	fmt.Printf("User %s has been registered\n", user.Name)

	return nil
}

func handlerReset(s *state, cmd command) error {
	if len(cmd.arguments) > 0 {
		return errors.New("Too many arguements: reset takes 0 arguments")
	}

	err := s.db.ResetUsers(context.Background())
	if err != nil {
		return err
	}
	fmt.Printf("Users have been reset\n")

	return nil
}

func handlerUsers(s *state, cmd command) error {
	if len (cmd.arguments) > 0 {
		return errors.New("Too many arguements: users takes 0 arguments")
	}
	users, err := s.db.GetUsers(context.Background())
	if err != nil {
		return err
	}
	for _, u := range users {
		if s.config.CurrentUserName == u.Name {
			fmt.Printf("* %s (current)\n", u.Name)
		} else {
			fmt.Printf("* %s\n", u.Name)
		}
	}

	return nil
}
