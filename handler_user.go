package main

import (
	"context"
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/BlessE1/gator/internal/database"
	"github.com/google/uuid"
)

func handlerLogin(s *state, cmd command) error {
	if len(cmd.Args) == 0 {
		return errors.New("Username is required")
	}

	ctx := context.Background()
	fetchedUser, err := s.db.GetUser(ctx, cmd.Args[0])
	if err != nil {
		fmt.Println("Username does not exist in database")
		os.Exit(1)
	}

	err = s.config.SetUser(cmd.Args[0])
	if err != nil {
		fmt.Println("Username could not be set in config")
		os.Exit(1)
	}

	fmt.Println("Username has been set in config")
	fmt.Println(fetchedUser)
	return nil
}

func handlerRegister(s *state, cmd command) error {
	if len(cmd.Args) == 0 {
		return errors.New("Username is required")
	}

	ctx := context.Background()
	createUserParams := database.CreateUserParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:      cmd.Args[0],
	}

	insertedUser, err := s.db.CreateUser(ctx, createUserParams)
	if err != nil {
		fmt.Println("Username could not be inserted")
		return err
	}

	err = handlerLogin(s, cmd)
	if err != nil {
		return err
	}

	fmt.Println("User has been created and logged in")
	fmt.Println(insertedUser)
	return nil
}

func handlerGetUsers(s *state, cmd command) error {
	ctx := context.Background()
	fetchedUsers, err := s.db.GetUsers(ctx)
	if err != nil {
		fmt.Println("Users could not be fetched")
		os.Exit(1)
	}

	for _, user := range fetchedUsers {
		if user.Name == s.config.CurrentUserName {
			fmt.Printf("* %s (current)\n", user.Name)
		} else {
			fmt.Printf("* %s\n", user.Name)
		}
	}
	return nil
}
