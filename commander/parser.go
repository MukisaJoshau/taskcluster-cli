package commander

import (
    "fmt"
)

func Parse(parent *Command, argv []string, start *int) (*Command, *CallContext) {
    fmt.Printf("Parsing args from step %d, parent is %s\n", *start, parent)
    fmt.Println("This function shall parse arguments, just gotta figure it out!")

    // start from the indicated (start) position in argv and parse the next command
    // and its arguments/options/flags, then return the new command as well as its context
    // if parent is nil it means we have the root, so don't bother about the command name

    return nil, nil
}
