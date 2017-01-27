package commander

import (
    "fmt"
    "os"
)

// the top of the command tree
// name of the command isn't important as we don't parse the first argv
var (
    root = NewCommand(os.Args[0])
)

// Root returns the root Command
func Root() *Command {
    return root
}

// Run starts the parser and the execution of sub*commands
func Run() bool {
    var command, final *Command
    var context *Context
    var err error
    step := 0

    for {
        // start by parsing
        final = command
        command, context, err = Parse(command, os.Args, &step)

        // on error, display usage
        if err != nil {
            fmt.Fprintf(os.Stderr, "error: %s\n", err)

            var usage string
            if command == nil {
                usage = Root().Usage()
            } else {
                usage = command.Usage()
            }
            fmt.Fprintf(os.Stderr, "%s", usage)

            return false

        // else stop parsing, we done here
        } else if command == nil && context == nil {
            break
        }

        // try to run Common from the CommandProvider of the parsed command
        if command.provider != nil {
            ok := command.provider.Common(context)
            if !ok {
                return false
            }
        }
    }

    // run Execute from the final command
    if final.provider != nil {
        ok := final.provider.Execute(context)
        if !ok {
            return false
        }
    }

    return true
}

/*
fmt.Printf("Command: %s\nFlags:\n", command.name)
for flag, _ := range context.Flags {
    fmt.Printf("\t%s\n", flag)
}
fmt.Printf("Options:\n")
for opt, val := range context.Options {
    fmt.Printf("\t%s: %s\n", opt, val)
}
fmt.Printf("Arguments:\n")
for arg, val := range context.Arguments {
    fmt.Printf("\t%s: %s\n", arg, val)
}
*/
