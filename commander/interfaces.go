package commander

// Context is the context given to a sub*command being executed
type Context struct {
	// Command line arguments and their values
	Arguments map[string]interface{}
	// Command line options and their values
	Options map[string]string
	// Command line flags
	Flags map[string]bool
}

func NewContext() *Context {
	return &Context{
		Arguments: make(map[string]interface{}),
		Options: make(map[string]string),
		Flags: make(map[string]bool),
	}
}

// CommandProvider describes the logic of a command and is passed to Command.Register()
type CommandProvider interface {
	// Common is called every time this command, or one of its subcommands, is requested
	// Return false to stop execution of further subcommands
	Common(context *Context) bool
	// Execute is called only when this command, but not its subcommands, is called
	Execute(context *Context) bool
}
