package script

type Config struct {
	Use                        string
	Aliases                    []string
	Short                      string
	Long                       string
	Example                    string
	ConfigPath                 string
	Hidden                     bool
	Annotations                map[string]string
	Version                    string
	suggestFor                 []string
	validArgs                  []string
	args                       PositionalArgs
	argAliases                 []string
	bashCompletionFunction     string
	deprecated                 string
	silenceErrors              bool
	silenceUsage               bool
	disableFlagParsing         bool
	disableAutoGenTag          bool
	disableFlagsInUseLine      bool
	disableSuggestions         bool
	suggestionsMinimumDistance int
	traverseChildren           bool
	fParseErrWhitelist         FParseErrWhitelist
}

func (c *Command) Initialize() {
	c.v.AddConfigPath(c.Config.ConfigPath)
}
