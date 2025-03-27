package node

import "fmt"

// CommandManager manages registered commands.
type CommandManager struct {
	commands map[string]Command
}

// NewCommandManager creates a new CommandManager.
func NewCommandManager() *CommandManager {
	return &CommandManager{
		commands: make(map[string]Command),
	}
}

// RegisterCommand registers a new command.
func (m *CommandManager) RegisterCommand(cmd Command) {
	m.commands[cmd.Name()] = cmd
}

// ExecuteCommand executes a command by name.
func (m *CommandManager) ExecuteCommand(name string) error {
	cmd, exists := m.commands[name]
	if !exists {
		return fmt.Errorf("command %s not found", name)
	}
	return cmd.Execute()
}
