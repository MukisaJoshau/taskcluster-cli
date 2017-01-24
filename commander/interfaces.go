package commander

// CallContext is the context given to a sub*command being executed
type CallContext struct {
	// Command line arguments and their values
	Arguments map[string]interface{}
	// Command line options and their values
	Options map[string]string
	// Command line flags
	Flags map[string]bool
}

// CommandLogic describes the logic of a command and is passed to Command.Register()
type CommandLogic interface {
	// Common is called every time this command, or one of its subcommands, is requested
	// Return false to stop execution of further subcommands
	Common(context CallContext) bool
	// Execute is called only when this command, but not its subcommands, is called
	Execute(context CallContext) bool
}
