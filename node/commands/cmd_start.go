package node

import "fmt"

// StartCommand is a concrete implementation of the Command interface.
type StartCommand struct{}

func (c *StartCommand) Execute() error {
	fmt.Println("Executing Start Command")
	// TODO

	return nil
}

func (c *StartCommand) Name() string {
	return "StartCommand"
}

// StopCommand is another implementation of the Command interface.
type StopCommand struct{}

func (c *StopCommand) Execute() error {
	fmt.Println("Executing Stop Command")
	return nil
}

func (c *StopCommand) Name() string {
	return "StopCommand"
}
