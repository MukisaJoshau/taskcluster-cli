package commander

import (
    "fmt"
    "regexp"
)

var longOptRegex *regexp.Regexp
var shortOptRegex *regexp.Regexp

func init() {
    var err error

    longOptRegex, err = regexp.Compile("^-(-[a-zA-Z0-9]+)+$")
    if(err != nil) {
        panic(fmt.Sprintf("longOptRegex could not be compiled: %s", err))
    }

    shortOptRegex, err = regexp.Compile("^-[a-zA-Z]+$")
    if(err != nil) {
        panic(fmt.Sprintf("shortOptRegex could not be compiled: %s", err))
    }
}

func Parse(parent *Command, argv []string, start *int) (*Command, *Context, error) {
    // if we're at the end of argv
    if len(argv) <= *start {
        return nil, nil, nil
    }

    // build args string
    args := argv[*start:]
    context := NewContext()

    // parse the command
    var cmd *Command
    if parent == nil {
        cmd = Root()
    } else {
        var ok bool
        cmd, ok = parent.subcommands[args[0]]
        if !ok {
            return parent, nil, fmt.Errorf("undefined subcommand '%s'", args[0])
        }
    }

    // now parse options
    var a int

    for a = 1; a < len(args); a++ {
        // long option
        if longOptRegex.MatchString(args[a]) {
            opt := args[a][2:]

            // option
            if _, ok := cmd.options[opt]; ok {
                a++
                if a < len(args) { // pick next term as option value
                    context.Options[opt] = args[a]
                } else {
                    return cmd, nil, fmt.Errorf("please supply value for option '--%s'", opt)
                }

            // flag
            } else if _, ok := cmd.flags[opt]; ok {
                context.Flags[opt] = true

            // not existing anywhere
            } else {
                return cmd, nil, fmt.Errorf("undefined option/flag '--%s' for command '%s'", opt, cmd.name)
            }

        // short option(s)
        } else if shortOptRegex.MatchString(args[a]) {
            // get all flags except the last, which might be an option
            var c int

            for c = 1; c < len(args[a]) - 1; c++ {
                char := args[a][c:c+1]
                if _, ok := cmd.flags[char]; ok {
                    context.Flags[char] = true
                } else {
                    return cmd, nil, fmt.Errorf("undefined flag '-%s'", char)
                }
            }

            // get last character, check for option or flag
            last := args[a][c:c+1]
            if _, ok := cmd.flags[last]; ok {
                context.Flags[last] = true

            // option
            } else if _, ok := cmd.options[last]; ok {
                a++
                if a < len(args) { // pick next term as option value
                    context.Options[last] = args[a]
                } else {
                    return cmd, nil, fmt.Errorf("please supply value for option '-%s'", last)
                }
            }

        // move on
        } else {
            break
        }
    }

    // arguments, if possible
    if len(cmd.subcommands) == 0 {
        var arg int

        for arg = 0; arg < len(cmd.arguments) && a < len(args); arg++ {
            argdesc := cmd.arguments[arg]

            // simple arguments, match them
            if argdesc.list == false {
                context.Arguments[argdesc.name] = args[a]

            // list argument
            } else {
                list := make([]string, 0)
                for ; a < len(args); a++ {
                    list = append(list, args[a])
                }
                context.Arguments[argdesc.name] = list
            }

            a++
        }

        // either we've run out of argv for all our ArgumentDescriptors
        // check if any remaning is required
        if arg < len(cmd.arguments) && cmd.arguments[arg].required {
            return cmd, nil, fmt.Errorf("missing required argument '<%s>'", cmd.arguments[arg].name)

        // or we've run out of ArgumentDescriptors to satisfy our argv
        // check for unused argv
        } else if a < len(args) {
            return cmd, nil, fmt.Errorf("too many arguments supplied for command '%s'", cmd.name)
        }
    }

    // at this point we should be pointing to the next subcommand or outside the boundaries of argv
    *start = *start + a

    return cmd, context, nil
}

/* TODO or just ideas
    - change name of root command
    - show usage for root command when error?
    - prevent typing same option twice
    - or enable typing same option multiple times
    - detect --help
    - prevent running with empty leaves
*/
