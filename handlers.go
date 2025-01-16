package main

import (
	"errors"
)

func handlerLogin(s *state, cmd command) error {
	if len(cmd.args) == 0 {
		return errors.New("expected a single argument, username")
	}

	s.config.CurrentUserName = cmd.args[0]
	s.config.SetUser(cmd.args[0])

	return nil
}
