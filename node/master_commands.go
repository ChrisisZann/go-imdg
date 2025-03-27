package node

import "fmt"

type StopCommand struct{}

func (c *StopCommand) Name() string {
	return "StopCommand"
}

type StartCommand struct{}

func (c *StartCommand) Name() string {
	return "StartCommand"
}

func (c *StartCommand) Execute() error {
	fmt.Println("StartCommand executed")
	return nil
}

func (c *StopCommand) Execute() error {
	fmt.Println("StopCommand executed")
	return nil
}

func (m *Master) initMasterCommands() {
	manager := NewCommandManager()

	// Register commands
	manager.RegisterCommand(&StartCommand{})
	manager.RegisterCommand(&StopCommand{})

	// // Execute commands dynamically
	// if err := manager.ExecuteCommand("StartCommand"); err != nil {
	// 	fmt.Println("Error:", err)
	// }

	// if err := manager.ExecuteCommand("StopCommand"); err != nil {
	// 	fmt.Println("Error:", err)
	// }
}
