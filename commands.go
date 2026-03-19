package main

import (
	"Tamelien/blog-aggregator/internal/config"
	"Tamelien/blog-aggregator/internal/database"
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
)

type state struct {
	db  *database.Queries
	cfg *config.Config
}

type command struct {
	name string
	args []string
}

type commands struct {
	handlerFunctions map[string]func(*state, command) error
}

func (c *commands) run(s *state, cmd command) error {
	handler, ok := c.handlerFunctions[cmd.name]
	if !ok {
		return fmt.Errorf("function not found")
	}

	return handler(s, cmd)
}

func (c *commands) register(name string, f func(*state, command) error) {
	if _, ok := c.handlerFunctions[name]; ok {
		return
	}

	c.handlerFunctions[name] = f
}

func handlerLogin(s *state, cmd command) error {
	if len(cmd.args) == 0 {
		return fmt.Errorf("usage: register <username>")
	}

	user, err := s.db.GetUser(context.Background(), cmd.args[0])
	if err != nil {
		return err
	}

	if err := s.cfg.SetUser(user.Name); err != nil {
		return err
	}

	fmt.Printf("User has been set to %s.\n", cmd.args[0])
	s.cfg.SetUser(user.Name)
	return nil
}

func handlerRegister(s *state, cmd command) error {
	if len(cmd.args) == 0 {
		return fmt.Errorf("usage: register <username>")
	}

	user, err := s.db.CreateUser(context.Background(),
		database.CreateUserParams{
			ID:        uuid.New(),
			Name:      cmd.args[0],
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		})
	if err != nil {
		return err
	}

	fmt.Printf("User %s has been created.\n", cmd.args[0])
	s.cfg.SetUser(user.Name)
	fmt.Printf("%+v\n", user)
	return nil
}
