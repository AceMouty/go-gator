package main

import (
	"context"
	"errors"
	"fmt"
	"github.com/acemouty/gator/internal/database"
)

func middlewareValidateUser(handler func(s *state, cmd command, user database.User) error) func(*state, command) error {
	return func(s *state, cmd command) error {
		user, err := s.db.GetUser(context.Background(), s.cfg.CurrentUserName)
		if err != nil {
			return errors.New("Unable to verify user")
		}

		return handler(s, cmd, user)
	}
}

func middlewareUserExists(handler func(s *state, cmd command) error) func(*state, command) error {
	return func(s *state, cmd command) error {
		userName := cmd.args[0]
		userExists, err := s.db.UserExists(context.Background(), userName)
		if err != nil {
			return errors.New("ran into a issue checking if user exists")
		}

		if userExists {
			errMsg := fmt.Sprintf("User with name %v already registered", userName)
			return errors.New(errMsg)
		}

		return handler(s, cmd)
	}
}
