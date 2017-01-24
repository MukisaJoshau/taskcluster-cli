package commander

import (
    "fmt"
    "strings"
)

// Command represents a node in the tree of commands
type Command struct {
    name        string
    description string
    parent      *Command
    subcommands map[string]*Command
    arguments   []*ArgumentDescriptor
    options     map[string]string // maps name -> description
    flags       map[string]string
    provider    *CommandLogic
}

// NewCommand creates a new Command object with the specified name
func NewCommand(name string) *Command {
    return &Command{
        name:        name,
        description: "no desc provided",
        parent:      nil,
        subcommands: make(map[string]*Command),
        arguments:   make([]*ArgumentDescriptor, 0),
        options:     make(map[string]string),
        flags:       make(map[string]string),
        provider:    nil,
    }
}

// ArgumentDescriptor represents an argument (<arg>) and its properties
type ArgumentDescriptor struct {
    name        string
    description string
    required    bool
    list        bool
}


// NewArgumentDescriptor creates a new ArgumentDescriptor
func NewArgumentDescriptor(name string, desc string, req bool, list bool) *ArgumentDescriptor {
    return &ArgumentDescriptor{
        name:        name,
        description: desc,
        required:    req,
        list:        list,
    }
}

// Get gets subcommand object, create one if none exists
func (this *Command) Get(name string) *Command {
    cmd, ok := this.subcommands[name];

    if !ok { // no command registered under that name
        if len(this.arguments) > 0 {
            panic("You cannot add subcommands to a command that has arguments")
        }

        cmd = NewCommand(name)
        cmd.parent = this

        this.subcommands[name] = cmd
    }

    return cmd
}

// Register registers a CommandLogic for the current Command
func (this *Command) Register(prov *CommandLogic) *Command {
    this.provider = prov
    return this
}

// Provider returns the CommandLogic associated to the current Command
func (this *Command) Provider() *CommandLogic {
    return this.provider
}

// Description sets the description for the current Command
func (this *Command) Description(desc string) *Command {
    this.description = desc
    return this
}

// Option adds option (--option VALUE) to the current Command
func (this *Command) Option(name string, desc string) *Command {
    this.options[name] = desc
    return this
}

// Flag adds a flag (like an option but no value) to the current Command
func (this *Command) Flag(name string, desc string) *Command {
    this.flags[name] = desc
    return this
}

// RequiredArgument adds an argument to the current Command
// You must register all required arguments before the optional ones
// Only the last argument can be flagged as a list argument
func (this *Command) Argument(name string, desc string, req bool, list bool) *Command {
    if len(this.subcommands) > 0 {
        panic("You cannot add arguments to a command that has subcommands.")
    }
    if len(this.arguments) > 0 {
        if req == true && this.arguments[len(this.arguments)-1].required == false {
            panic("You cannot add required arguments to a command that already has optional arguments.")
        }
        if this.arguments[len(this.arguments)-1].list == true {
            panic("You cannot add any more arguments after adding a list argument.")
        }
    }

    this.arguments = append(this.arguments, NewArgumentDescriptor(name, desc, req, list))

    return this
}

// Usage generates a usage string for this Command
func (this *Command) Usage() string {
    usage := this.ShortUsage()

    // display arguments after usage string
    for _, arg := range this.arguments {
        argstr := fmt.Sprintf("<%s>", arg.name)
        if arg.list {
            argstr += "..."
        }
        if !arg.required {
            argstr = fmt.Sprintf("[%s]", argstr)
        }

        usage += " " + argstr
    }

    // display arguments and their descriptions
    if len(this.arguments) > 0 {
        usage += "\n\nArguments:"
        for _, arg := range this.arguments {
            usage += fmt.Sprintf("\n\t<%s>: %s", arg.name, arg.description)
        }
    }

    // display subcommands and their descriptions
    if len(this.subcommands) > 0 {
        usage += "\n\nSubcommands:"
        for name, cmd := range this.subcommands {
            usage += fmt.Sprintf("\n\t%s: %s", name, cmd.description)
        }
    }

    // display options and their descriptions
    if len(this.options) > 0 {
        usage += "\n\nOptions:"
        for name, desc := range this.options {
            usage += "\n\t -"
            if len(name) > 1 {
                usage += "-"
            }
            usage += fmt.Sprintf("%s %s: %s", name, strings.ToUpper(name), desc)
        }
    }

    // display flags and their descriptions
    if len(this.flags) > 0 {
        usage += "\n\nFlags:"
        for name, desc := range this.flags {
            usage += "\n\t -"
            if len(name) > 1 {
                usage += "-"
            }
            usage += fmt.Sprintf("%s: %s", name, desc)
        }
    }

    usage += "\n"
    return usage
}

// ShortUsage generates a short usage string for this Command
func (this *Command) ShortUsage() string {
    var usage string

    // we get stuff recursively
    if this.parent == nil {
        usage += "Usage:"
    } else {
        usage += this.parent.ShortUsage()
    }

    usage += fmt.Sprintf(" %s", this.name)
    return usage
}
