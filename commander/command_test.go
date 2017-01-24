package commander

import(
    "testing"

    assert "github.com/stretchr/testify/require"
)

/* TEST IDEAS

if I Get() a non-existent command, does it create it?
are descriptions updated properly?
are subcommnds, arguments, options and flags properly added?
can i add an optional argument before a required argument?
can i add an argument after a list argument?
can i add arguments to a command that has subcommands?
can i add subcommands to a command that has arguments?
does Usage() print properly?
does ShortUsage() print properly?

*/

func TestGetNonExistingCommand(t *testing.T) {
    assert := assert.New(t)

    root := extpoints.Root()
    assert.Equal(len(root.subcommands), 0, "Root should start with no subcommands")

    root.Get("test")
    cmd, ok := root.subcommands["test"]
    assert.True(ok, "Get(\"test\") should create a subcommand attached to 'test' on the root Command")
    assert.Equal(cmd.name, "test", "Get(\"test\") should create a subcommand with 'test' as its name attribute")
}

func TestGetExistingCommand(t *testing.T) {
    assert := assert.New(t)

    root := extpoints.Root()
    root
}
