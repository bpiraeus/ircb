// This file is subject to a 1-clause BSD license.
// Its contents can be found in the enclosed LICENSE file.

package cmd

import (
	"github.com/jteeuwen/ircb/proto"
)

// CommandFunc represents a command constructor.
type CommandFunc func() *Command

// ExecuteFunc represents a command execution handler.
// These are executed in a separate goroutine.
type ExecuteFunc func(*Command, *proto.Client, *proto.Message)

// Command represents a single bot command.
type Command struct {
	Name        string      // Command name.
	Description string      // Command description.
	Params      []Param     // Command parameters.
	Execute     ExecuteFunc // Execution handler for the command.
	Restricted  bool        // Command is restricted to admin users.
}

// RequiredParamCount counts the number of required parameters.
func (c *Command) RequiredParamCount() int {
	var pc int

	for i := range c.Params {
		if !c.Params[i].Optional {
			pc++
		}
	}

	return pc
}

// List of registered command constructors.
var commands = make(map[string]CommandFunc)

// Register registers the given command name and constructor.
// Modules should call this during initialization to register their
// commands with the bot.
func Register(name string, cf CommandFunc) {
	commands[name] = cf
}
