package extpoints

import (
    "fmt"
    "strings"
)

// the top of the command tree
var (
    root = newCommand("taskcluster")
)

func Root() *Command {
    return root
}

// extpoints.Root().Register(taskcluster{}) <- for example

type Command struct {
    name        string
    description string
    parent      *Command
    subcommands map[string]*Command
    arguments   []string
    argDescs    []string
    argReqs     int
    options     map[string]string // maps name -> description
    flags       map[string]string
    provider    *CommandProvider
}

/*
 * Get subcommand object, create one if none exists
 */
func (this *Command) Get(name string) *Command {
    cmd, ok := this.subcommands[name];

    if !ok { // no command registered under that name
        if len(this.arguments) > 0 {
            panic("You cannot add subcommands to a command that has arguments")
        }

        cmd = newCommand(name)
        cmd.parent = this

        this.subcommands[name] = cmd
    }

    return cmd
}

/*
 * Create a new Command object with the specified name
 */
func newCommand(name string) *Command {
    cmd := &Command{
        name:        name,
        description: "no desc provided",
        parent:      nil,
        subcommands: make(map[string]*Command),
        arguments:   make([]string, 0),
        argDescs:    make([]string, 0),
        argReqs:     0,
        options:     make(map[string]string),
        flags:       make(map[string]string),
        provider:    nil,
    }

    return cmd
}

/*
 * Register a CommandProvider for the current Command
 */
func (this *Command) Register(prov *CommandProvider) *Command {
    this.provider = prov
    return this
}

/*
 * Get the CommandProvider associated to the current Command
 */
func (this *Command) Provider() *CommandProvider {
    return this.provider
}

/*
 * Set description for the current Command
 */
func (this *Command) Description(desc string) *Command {
    this.description = desc
    return this
}

/*
 * Add option (--option VALUE) to the current Command
 */
func (this *Command) Option(name string, desc string) *Command {
    this.options[name] = desc
    return this
}

/*
 * Add flag (like an option but no value) to the current Command
 */
func (this *Command) Flag(name string, desc string) *Command {
    this.flags[name] = desc
    return this
}

/*
 * Add a required argument to the current Command
 * You must register all required arguments before the optional ones
 */
func (this *Command) RequiredArgument(name string, desc string) *Command {
    if len(this.subcommands) > 0 {
        panic("You cannot add arguments to a command that has subcommands.")
    }
    if this.argReqs < len(this.arguments) { // means optional args have already been added
        panic("You cannot add required arguments to a command that already has optional arguments.")
    }

    this.arguments = append(this.arguments, name)
    this.argDescs = append(this.argDescs, desc)

    this.argReqs += 1

    return this
}

/*
 * Add an optional argument to the current Command
 * You must register all required arguments before calling this function
 */
func (this *Command) OptionalArgument(name string, desc string) *Command {
    if len(this.subcommands) > 0 {
        panic("You cannot add arguments to a command that has subcommands.")
    }

    this.arguments = append(this.arguments, name)
    this.argDescs = append(this.argDescs, desc)

    return this
}

/*
 * Generate a usage string for this Command
 */
func (this *Command) Usage() string {
    usage := this.ShortUsage()

    // display arguments after usage string
    i := 1
    for _, name := range this.arguments {
        if i <= this.argReqs { // still in required arguments
            usage += fmt.Sprintf(" <%s>", name)
        } else {
            usage += fmt.Sprintf(" [<%s>]", name)
        }
        i += 1
    }

    // display arguments and their descriptions
    if len(this.arguments) > 0 {
        usage += "\n\nArguments:"
        for pos, name := range this.arguments {
            usage += fmt.Sprintf("\n\t<%s>: %s", name, this.argDescs[pos])
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
            usage += fmt.Sprintf("\n\t%s %s: %s", name, strings.ToUpper(name), desc)
        }
    }

    // display flags and their descriptions
    if len(this.flags) > 0 {
        usage += "\n\nFlags:"
        for name, desc := range this.flags {
            usage += fmt.Sprintf("\n\t%s: %s", name, desc)
        }
    }

    usage += "\n"
    return usage
}

/*
 * Generate a short usage string for this Command
 */
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
