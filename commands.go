package main

import "errors"

type commandMap = map[string]func(*state, command) error

type commandStore struct {
	commandsMap commandMap
}

func (c *commandStore) register(commandName string, f func(*state, command) error) {
	c.commandsMap[commandName] = f
}

func (c *commandStore) run(s *state, cmd command) error {
	handlerCommand, ok := c.commandsMap[cmd.name]

	if !ok {
		return errors.New("Command not found")
	}

	return handlerCommand(s, cmd)
}
