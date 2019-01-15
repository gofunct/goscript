package script

import (
	"bytes"
	flag "github.com/spf13/pflag"
	"io"
)

var initializers []func()
var EnablePrefixMatching = false
var EnableCommandSorting = true
var minNamePadding = 11
var minCommandPathPadding = 11
var minUsagePadding = 25

type FParseErrWhitelist flag.ParseErrorsWhitelist

type Command struct {
	Use                        string
	Aliases                    []string
	SuggestFor                 []string
	Short                      string
	Long                       string
	Example                    string
	ValidArgs                  []string
	Args                       PositionalArgs
	ArgAliases                 []string
	BashCompletionFunction     string
	Deprecated                 string
	Hidden                     bool
	Annotations                map[string]string
	Version                    string
	PersistentPreRun           func(cmd *Command, args []string)
	PersistentPreRunE          func(cmd *Command, args []string) error
	PreRun                     func(cmd *Command, args []string)
	PreRunE                    func(cmd *Command, args []string) error
	Run                        func(cmd *Command, args []string)
	RunE                       func(cmd *Command, args []string) error
	PostRun                    func(cmd *Command, args []string)
	PostRunE                   func(cmd *Command, args []string) error
	PersistentPostRun          func(cmd *Command, args []string)
	PersistentPostRunE         func(cmd *Command, args []string) error
	SilenceErrors              bool
	SilenceUsage               bool
	DisableFlagParsing         bool
	DisableAutoGenTag          bool
	DisableFlagsInUseLine      bool
	DisableSuggestions         bool
	SuggestionsMinimumDistance int
	TraverseChildren           bool
	FParseErrWhitelist         FParseErrWhitelist
	commands                   []*Command
	parent                     *Command
	commandsMaxUseLen          int
	commandsMaxCommandPathLen  int
	commandsMaxNameLen         int
	commandsAreSorted          bool
	commandCalledAs            struct {
		name   string
		called bool
	}
	args            []string
	flagErrorBuf    *bytes.Buffer
	flags           *flag.FlagSet
	pflags          *flag.FlagSet
	lflags          *flag.FlagSet
	iflags          *flag.FlagSet
	parentsPflags   *flag.FlagSet
	globNormFunc    func(f *flag.FlagSet, name string) flag.NormalizedName
	output          io.Writer
	usageFunc       func(*Command) error
	usageTemplate   string
	flagErrorFunc   func(*Command, error) error
	helpTemplate    string
	helpFunc        func(*Command, []string)
	helpCommand     *Command
	versionTemplate string
}
