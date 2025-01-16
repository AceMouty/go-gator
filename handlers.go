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

	s.cfg.CurrentUserName = cmd.args[0]
	s.cfg.SetUser(cmd.args[0])

	return nil
}

func handlerRegitser(s *state, cmd command) error {
	if len(cmd.args) == 0 {
		return errors.New("expected a username to be provided")
	}

	userName := cmd.args[0]
	userExists, err := s.db.UserExists(context.Background(), userName)
	if err != nil {
		return err
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
