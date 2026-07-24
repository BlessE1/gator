package main

import (
 "errors"
)

type command struct {
 Name string
 Args []string
}

type commands struct {
 cmdMap map[string]func(*state,  command) error
}

func (c *commands) run(s *state, cmd command) error {
 handler, exists := c.cmdMap[cmd.Name]
 if exists {
  return handler(s, cmd)
 } else {
  return errors.New("command does not exist")
 }
}

func (c *commands) register(name string, f func(*state, command) error) { 
 c.cmdMap[name] = f
}
