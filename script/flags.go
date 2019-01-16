package script

import (
	"bytes"
	"fmt"
	"github.com/gofunct/goscript/utils"
	"github.com/spf13/pflag"
	"strings"
)

const OneRequiredFlag = "goscript_annotation_one_required_flag"

// GlobalNormalizationFunc returns the global normalization function or nil if it doesn't exist.
func (c *Command) GlobalNormalizationFunc() func(f *pflag.FlagSet, name string) pflag.NormalizedName {
	return c.globNormFunc
}

// Flags returns the complete FlagSet that applies
// to this command (local and persistent declared here and by all parents).
func (c *Command) Flags() *pflag.FlagSet {
	if c.flags == nil {
		c.flags = pflag.NewFlagSet(c.Name(), pflag.ContinueOnError)
		if c.flagErrorBuf == nil {
			c.flagErrorBuf = new(bytes.Buffer)
		}
		c.flags.SetOutput(c.flagErrorBuf)
	}

	return c.flags
}

// LocalNonPersistentFlags are flags specific to this command which will NOT persist to subcommands.
func (c *Command) LocalNonPersistentFlags() *pflag.FlagSet {
	persistentFlags := c.PersistentFlags()

	out := pflag.NewFlagSet(c.Name(), pflag.ContinueOnError)
	c.LocalFlags().VisitAll(func(f *pflag.Flag) {
		if persistentFlags.Lookup(f.Name) == nil {
			out.AddFlag(f)
		}
	})
	return out
}

// LocalFlags returns the local FlagSet specifically set in the current command.
func (c *Command) LocalFlags() *pflag.FlagSet {
	c.mergePersistentFlags()

	if c.lflags == nil {
		c.lflags = pflag.NewFlagSet(c.Name(), pflag.ContinueOnError)
		if c.flagErrorBuf == nil {
			c.flagErrorBuf = new(bytes.Buffer)
		}
		c.lflags.SetOutput(c.flagErrorBuf)
	}
	c.lflags.SortFlags = c.Flags().SortFlags
	if c.globNormFunc != nil {
		c.lflags.SetNormalizeFunc(c.globNormFunc)
	}

	addToLocal := func(f *pflag.Flag) {
		if c.lflags.Lookup(f.Name) == nil && c.parentsPflags.Lookup(f.Name) == nil {
			c.lflags.AddFlag(f)
		}
	}
	c.Flags().VisitAll(addToLocal)
	c.PersistentFlags().VisitAll(addToLocal)
	return c.lflags
}

// InheritedFlags returns all flags which were inherited from parents commands.
func (c *Command) InheritedFlags() *pflag.FlagSet {
	c.mergePersistentFlags()

	if c.iflags == nil {
		c.iflags = pflag.NewFlagSet(c.Name(), pflag.ContinueOnError)
		if c.flagErrorBuf == nil {
			c.flagErrorBuf = new(bytes.Buffer)
		}
		c.iflags.SetOutput(c.flagErrorBuf)
	}

	local := c.LocalFlags()
	if c.globNormFunc != nil {
		c.iflags.SetNormalizeFunc(c.globNormFunc)
	}

	c.parentsPflags.VisitAll(func(f *pflag.Flag) {
		if c.iflags.Lookup(f.Name) == nil && local.Lookup(f.Name) == nil {
			c.iflags.AddFlag(f)
		}
	})
	return c.iflags
}

// NonInheritedFlags returns all flags which were not inherited from parent commands.
func (c *Command) NonInheritedFlags() *pflag.FlagSet {
	return c.LocalFlags()
}

// PersistentFlags returns the persistent FlagSet specifically set in the current command.
func (c *Command) PersistentFlags() *pflag.FlagSet {
	if c.pflags == nil {
		c.pflags = pflag.NewFlagSet(c.Name(), pflag.ContinueOnError)
		if c.flagErrorBuf == nil {
			c.flagErrorBuf = new(bytes.Buffer)
		}
		c.pflags.SetOutput(c.flagErrorBuf)
	}
	return c.pflags
}

// ResetFlags deletes all flags from command.
func (c *Command) ResetFlags() {
	c.flagErrorBuf = new(bytes.Buffer)
	c.flagErrorBuf.Reset()
	c.flags = pflag.NewFlagSet(c.Name(), pflag.ContinueOnError)
	c.flags.SetOutput(c.flagErrorBuf)
	c.pflags = pflag.NewFlagSet(c.Name(), pflag.ContinueOnError)
	c.pflags.SetOutput(c.flagErrorBuf)

	c.lflags = nil
	c.iflags = nil
	c.parentsPflags = nil
}

// HasFlags checks if the command contains any flags (local plus persistent from the entire structure).
func (c *Command) HasFlags() bool {
	return c.Flags().HasFlags()
}

// HasPersistentFlags checks if the command contains persistent flags.
func (c *Command) HasPersistentFlags() bool {
	return c.PersistentFlags().HasFlags()
}

// HasLocalFlags checks if the command has flags specifically declared locally.
func (c *Command) HasLocalFlags() bool {
	return c.LocalFlags().HasFlags()
}

// HasInheritedFlags checks if the command has flags inherited from its parent command.
func (c *Command) HasInheritedFlags() bool {
	return c.InheritedFlags().HasFlags()
}

// HasAvailableFlags checks if the command contains any flags (local plus persistent from the entire
// structure) which are not hidden or deprecated.
func (c *Command) HasAvailableFlags() bool {
	return c.Flags().HasAvailableFlags()
}

// HasAvailablePersistentFlags checks if the command contains persistent flags which are not hidden or deprecated.
func (c *Command) HasAvailablePersistentFlags() bool {
	return c.PersistentFlags().HasAvailableFlags()
}

// HasAvailableLocalFlags checks if the command has flags specifically declared locally which are not hidden
// or deprecated.
func (c *Command) HasAvailableLocalFlags() bool {
	return c.LocalFlags().HasAvailableFlags()
}

// HasAvailableInheritedFlags checks if the command has flags inherited from its parent command which are
// not hidden or deprecated.
func (c *Command) HasAvailableInheritedFlags() bool {
	return c.InheritedFlags().HasAvailableFlags()
}

// Flag climbs up the command tree looking for matching flag.
func (c *Command) Flag(name string) (flag *pflag.Flag) {
	flag = c.Flags().Lookup(name)

	if flag == nil {
		flag = c.persistentFlag(name)
	}

	return
}

// Recursively find matching persistent flag.
func (c *Command) persistentFlag(name string) (flag *pflag.Flag) {
	if c.HasPersistentFlags() {
		flag = c.PersistentFlags().Lookup(name)
	}

	if flag == nil {
		c.updateParentsPflags()
		flag = c.parentsPflags.Lookup(name)
	}
	return
}

// ParseFlags parses persistent flag tree and local flags.
func (c *Command) ParseFlags(args []string) error {
	if c.Config.disableFlagParsing {
		return nil
	}

	if c.flagErrorBuf == nil {
		c.flagErrorBuf = new(bytes.Buffer)
	}
	beforeErrorBufLen := c.flagErrorBuf.Len()
	c.mergePersistentFlags()

	//do it here after merging all flags and just before parse
	c.Flags().ParseErrorsWhitelist = pflag.ParseErrorsWhitelist(c.Config.fParseErrWhitelist)

	err := c.Flags().Parse(args)
	// Print warnings if they occurred (e.g. deprecated flag messages).
	if c.flagErrorBuf.Len()-beforeErrorBufLen > 0 && err == nil {
		c.Debug(c.flagErrorBuf.String())
	}

	return err
}

// Parent returns a commands parent command.
func (c *Command) Parent() *Command {
	return c.parent
}

// mergePersistentFlags merges c.PersistentFlags() to c.Flags()
// and adds missing persistent flags of all parents.
func (c *Command) mergePersistentFlags() {
	c.updateParentsPflags()
	c.Flags().AddFlagSet(c.PersistentFlags())
	c.Flags().AddFlagSet(c.parentsPflags)
}

// updateParentsPflags updates c.parentsPflags by adding
// new persistent flags of all parents.
// If c.parentsPflags == nil, it makes new.
func (c *Command) updateParentsPflags() {
	if c.parentsPflags == nil {
		c.parentsPflags = pflag.NewFlagSet(c.Name(), pflag.ContinueOnError)
		c.parentsPflags.SetOutput(c.flagErrorBuf)
		c.parentsPflags.SortFlags = false
	}

	if c.globNormFunc != nil {
		c.parentsPflags.SetNormalizeFunc(c.globNormFunc)
	}

	c.Root().PersistentFlags().AddFlagSet(pflag.CommandLine)

	c.VisitParents(func(parent *Command) {
		c.parentsPflags.AddFlagSet(parent.PersistentFlags())
	})
}

// DebugFlags used to determine which flags have been assigned to which commands
// and which persist.
func (c *Command) DebugFlags() {
	c.Debug("DebugFlags called on", c.Name())
	var debugflags func(*Command)

	debugflags = func(x *Command) {
		if x.HasFlags() || x.HasPersistentFlags() {
			c.Debug(x.Name())
		}
		if x.HasFlags() {
			x.flags.VisitAll(func(f *pflag.Flag) {
				if x.HasPersistentFlags() && x.persistentFlag(f.Name) != nil {
					c.Debug("  -"+f.Shorthand+",", "--"+f.Name, "["+f.DefValue+"]", "", f.Value.String(), "  [LP]")
				} else {
					c.Debug("  -"+f.Shorthand+",", "--"+f.Name, "["+f.DefValue+"]", "", f.Value.String(), "  [L]")
				}
			})
		}
		if x.HasPersistentFlags() {
			x.pflags.VisitAll(func(f *pflag.Flag) {
				if x.HasFlags() {
					if x.flags.Lookup(f.Name) == nil {
						c.Debug("  -"+f.Shorthand+",", "--"+f.Name, "["+f.DefValue+"]", "", f.Value.String(), "  [P]")
					}
				} else {
					c.Debug("  -"+f.Shorthand+",", "--"+f.Name, "["+f.DefValue+"]", "", f.Value.String(), "  [P]")
				}
			})
		}
		c.Debug(x.flagErrorBuf.String())
		if x.HasSubCommands() {
			for _, y := range x.commands {
				debugflags(y)
			}
		}
	}

	debugflags(c)
}

func stripFlags(args []string, c *Command) []string {
	if len(args) == 0 {
		return args
	}
	c.mergePersistentFlags()

	commands := []string{}
	flags := c.Flags()

Loop:
	for len(args) > 0 {
		s := args[0]
		args = args[1:]
		switch {
		case s == "--":
			// "--" terminates the flags
			break Loop
		case strings.HasPrefix(s, "--") && !strings.Contains(s, "=") && !utils.HasNoOptDefVal(s[2:], flags):
			// If '--flag arg' then
			// delete arg from args.
			fallthrough // (do the same as below)
		case strings.HasPrefix(s, "-") && !strings.Contains(s, "=") && len(s) == 2 && !utils.ShortHasNoOptDefVal(s[1:], flags):
			// If '-f arg' then
			// delete 'arg' from args or break the loop if len(args) <= 1.
			if len(args) <= 1 {
				break Loop
			} else {
				args = args[1:]
				continue
			}
		case s != "" && !strings.HasPrefix(s, "-"):
			commands = append(commands, s)
		}
	}

	return commands
}

// MarkFlagRequired adds the BashCompOneRequiredFlag annotation to the named flag if it exists,
// and causes your command to report an error if invoked without the flag.
func (c *Command) MarkFlagRequired(name string) error {
	return MarkFlagRequired(c.Flags(), name)
}

// MarkPersistentFlagRequired adds the BashCompOneRequiredFlag annotation to the named persistent flag if it exists,
// and causes your command to report an error if invoked without the flag.
func (c *Command) MarkPersistentFlagRequired(name string) error {
	return MarkFlagRequired(c.PersistentFlags(), name)
}

// MarkFlagRequired adds the BashCompOneRequiredFlag annotation to the named flag if it exists,
// and causes your command to report an error if invoked without the flag.
func MarkFlagRequired(flags *pflag.FlagSet, name string) error {
	return flags.SetAnnotation(name, OneRequiredFlag, []string{"true"})
}
func (c *Command) validateRequiredFlags() error {
	flags := c.Flags()
	missingFlagNames := []string{}
	flags.VisitAll(func(pflag *pflag.Flag) {
		requiredAnnotation, found := pflag.Annotations[OneRequiredFlag]
		if !found {
			return
		}
		if (requiredAnnotation[0] == "true") && !pflag.Changed {
			missingFlagNames = append(missingFlagNames, pflag.Name)
		}
	})

	if len(missingFlagNames) > 0 {
		return fmt.Errorf(`required flag(s) "%s" not set`, strings.Join(missingFlagNames, `", "`))
	}
	return nil
}

// SetFlagErrorFunc sets a function to generate an error when flag parsing
// fails.
func (c *Command) SetFlagErrorFunc(f func(*Command, error) error) {
	c.flagErrorFunc = f
}

// FlagErrorFunc returns either the function set by SetFlagErrorFunc for this
// command or a parent, or it returns a function which returns the original
// error.
func (c *Command) FlagErrorFunc() (f func(*Command, error) error) {
	if c.flagErrorFunc != nil {
		return c.flagErrorFunc
	}

	if c.HasParent() {
		return c.parent.FlagErrorFunc()
	}
	return func(c *Command, err error) error {
		return err
	}
}
