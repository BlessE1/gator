package main

import (
	"context"
	"fmt"
	"os"
)

func handlerReset(s *state, cmd command) error {
	ctx := context.Background()
	err := s.db.ResetDatabase(ctx)
	if err != nil {
		fmt.Println("Database could not be reset")
		os.Exit(1)
	}

	fmt.Println("Database has been reset")
	return nil
}
