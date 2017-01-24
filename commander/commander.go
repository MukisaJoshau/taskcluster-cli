package commander

import (
    "fmt"
    "os"
)

// the top of the command tree
// name of the command isn't important as we don't parse the first argv
var (
    root = NewCommand("root")
)

// Root returns the root Command
func Root() *Command {
    return root
}

// Run starts the parser and the execution of sub*commands
func Run() bool {
    var command *Command
    var step int

    for {
        command, context := Parse(command, os.Args, &step)
        if command == nil && context == nil {
            break
        }

        // other things
    }

    fmt.Printf("Soon we shall be running!\n")

    return true
}
