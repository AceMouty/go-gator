package main

import (
	"context"
	"errors"
	"github.com/acemouty/gator/internal/database"
	"github.com/google/uuid"
	"log"
	"time"
)

func handlerLogin(s *state, cmd command) error {
	if len(cmd.args) == 0 {
		return errors.New("expected a single argument, username")
	}

	userName := cmd.args[0]
	userExists, err := s.db.UserExists(context.Background(), userName)
	if err != nil {
		log.Fatalf("ran into a issue checking if user exists: %v", err)
	}

	if !userExists {
		log.Fatalf("user %v doesnt exist", userName)
	}

	s.cfg.CurrentUserName = userName
	s.cfg.SetUser(userName)

	return nil
}

func handlerRegitser(s *state, cmd command) error {
	if len(cmd.args) == 0 {
		return errors.New("expected a username to be provided")
	}

	userName := cmd.args[0]
	userExists, err := s.db.UserExists(context.Background(), userName)
	if err != nil {
		log.Fatalf("ran into a issue checking if user exists: %v", err)
	}

	if userExists {
		log.Fatalf("User with name %v already registered", userName)
	}

	newId := uuid.New()
	now := time.Now()
	queryData := database.CreateUserParams{
		ID:        newId,
		CreatedAt: now,
		UpdatedAt: now,
		Name:      userName,
	}

	_, err = s.db.CreateUser(context.Background(), queryData)
	if err != nil {
		log.Fatalf("unable to register user %v: reason; %v", userName, err)
	}

	log.Println("registered user successfully")
	s.cfg.SetUser(userName)

	return nil
}

func handlerReset(s *state, cmd command) error {
	err := s.db.DeleteUsers(context.Background())

	if err != nil {
		log.Fatalf("ran into a issue deleting users: %v", err)
	}

	return nil
}
