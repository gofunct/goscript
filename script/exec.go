package script

import (
	"flag"
	"fmt"
	kitlog "github.com/go-kit/kit/log"
	"github.com/mitchellh/go-homedir"
	"github.com/pkg/errors"
	"github.com/spf13/viper"
	"log"
	"os"
	"path/filepath"
)

func (c *Command) execute(a []string) (err error) {

	if c == nil {
		return errors.New("Called Execute() on a nil Command")
	}

	if len(c.Config.deprecated) > 0 {
		c.Debug("Command %q is deprecated, %s\n" + c.Name() + c.Config.deprecated)
	}

	{
		logger := kitlog.NewJSONLogger(kitlog.NewSyncWriter(c.OutOrStderr()))
		logger = kitlog.With(logger, "time", kitlog.DefaultTimestampUTC, "origin", kitlog.DefaultCaller, "cmd", c.Config.Use)
		log.SetOutput(kitlog.NewStdlibAdapter(logger))
		log.SetOutput(c.OutOrStderr())
	}

	{
		c.InitDefaultHelpFlag()
		c.InitDefaultVersionFlag()
	}

	err = c.ParseFlags(a)
	if err != nil {
		return c.FlagErrorFunc()(c, err)
	}

	// If help is called, regardless of other flags, return we want help.
	// Also say we need help if the command isn't runnable.
	helpVal, err := c.Flags().GetBool("help")
	if err != nil {
		// should be impossible to get here as we always declare a help
		// flag in InitDefaultHelpFlag()
		c.Debug("\"help\" flag declared as non-bool. Please correct your code")
		return err
	}

	if helpVal {
		return flag.ErrHelp
	}

	// for back-compat, only add version flag behavior if version is defined
	if c.Config.Version != "" {
		versionVal, err := c.Flags().GetBool("version")
		if err != nil {
			c.Debug("\"version\" flag declared as non-bool. Please correct your code")
			return err
		}
		if versionVal {
			err := Compile(c.OutOrStdout(), c.VersionTemplate(), c)
			if err != nil {
				c.Debug(err.Error())
			}
			return err
		}
	}

	if !c.Runnable() {
		return flag.ErrHelp
	}

	c.preRun()

	argWoFlags := c.Flags().Args()
	if c.Config.disableFlagParsing {
		argWoFlags = a
	}

	if err := c.ValidateArgs(argWoFlags); err != nil {
		return err
	}

	for p := c; p != nil; p = p.Parent() {
		if p.PersistentPreRunE != nil {
			if err := p.PersistentPreRunE(c, argWoFlags); err != nil {
				return err
			}
			break
		} else if p.PersistentPreRun != nil {
			p.PersistentPreRun(c, argWoFlags)
			break
		}
	}
	if c.PreRunE != nil {
		if err := c.PreRunE(c, argWoFlags); err != nil {
			return err
		}
	} else if c.PreRun != nil {
		c.PreRun(c, argWoFlags)
	}

	if err := c.validateRequiredFlags(); err != nil {
		return errors.Wrap(err, "failed to validate required flags")
	}
	if c.RunE != nil {
		if err := c.RunE(c, argWoFlags); err != nil {
			return err
		}
	} else {
		c.Run(c, argWoFlags)
	}
	if c.PostRunE != nil {
		if err := c.PostRunE(c, argWoFlags); err != nil {
			return err
		}
	} else if c.PostRun != nil {
		c.PostRun(c, argWoFlags)
	}
	for p := c; p != nil; p = p.Parent() {
		if p.PersistentPostRunE != nil {
			if err := p.PersistentPostRunE(c, argWoFlags); err != nil {
				return err
			}
			break
		} else if p.PersistentPostRun != nil {
			p.PersistentPostRun(c, argWoFlags)
			break
		}
	}

	return nil
}

func (c *Command) preRun() {
	for _, x := range initializers {
		x()
	}
}

// Execute uses the args (os.Args[1:] by default)
// and run through the command tree finding appropriate matches
// for commands and then corresponding flags.
func (c *Command) Execute() error {
	_, err := c.ExecuteC()
	return err
}

// ExecuteC executes the command.
func (c *Command) ExecuteC() (cmd *Command, err error) {
	// Regardless of what command execute is called on, run on Root only
	if c.HasParent() {
		return c.Root().ExecuteC()
	}

	// windows hook
	if preExecHookFn != nil {
		preExecHookFn(c)
	}

	// initialize help as the last point possible to allow for user
	// overriding
	c.InitDefaultHelpCmd()

	var args []string

	// Workaround FAIL with "go test -v" or "cobra.test -test.v", see #155
	if c.args == nil && filepath.Base(os.Args[0]) != "cobra.test" {
		args = os.Args[1:]
	} else {
		args = c.args
	}

	var flags []string
	if c.Config.traverseChildren {
		cmd, flags, err = c.Traverse(args)
	} else {
		cmd, flags, err = c.Find(args)
	}
	if err != nil {
		// If found parse to a subcommand and then failed, talk about the subcommand
		if cmd != nil {
			c = cmd
		}
		if !c.Config.silenceErrors {
			c.Debug("Error:\n" + err.Error())
			c.Debug("Run '%v --help' for usage.\n" + c.CommandPath())
		}
		return c, err
	}

	cmd.commandCalledAs.called = true
	if cmd.commandCalledAs.name == "" {
		cmd.commandCalledAs.name = cmd.Name()
	}

	err = cmd.execute(flags)
	if err != nil {
		// Always show help if requested, even if SilenceErrors is in
		// effect
		if err == flag.ErrHelp {
			cmd.HelpFunc()(cmd, args)
			return cmd, nil
		}

		// If root command has SilentErrors flagged,
		// all subcommands should respect it
		if !cmd.Config.silenceErrors && !c.Config.silenceErrors {
			c.Debug("Error:\n" + err.Error())
		}

		// If root command has SilentUsage flagged,
		// all subcommands should respect it
		if !cmd.Config.silenceUsage && !c.Config.silenceUsage {
			c.Debug(cmd.UsageString())
		}
	}
	return cmd, err
}

// OnInitialize sets the passed functions to be run when each command's
// Execute method is called.
func OnInitialize(y ...func()) {
	initializers = append(initializers, y...)
}

func InitConfig(cfg string) func() {
	return func() {
		logger := kitlog.NewJSONLogger(kitlog.NewSyncWriter(os.Stdout))
		logger = kitlog.With(logger, "ts", kitlog.DefaultTimestampUTC, "caller", kitlog.DefaultCaller, "user", os.Getenv("USER"))
		log.SetOutput(kitlog.NewStdlibAdapter(logger))

		if cfg != "" {
			// Use config file from the flag.
			viper.SetConfigFile(cfg)
		} else {
			// Find home directory.
			home, err := homedir.Dir()
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}

			// Search config in home directory with name ".goscript" (without extension).
			viper.AddConfigPath(home)
			viper.SetConfigName(".chronic")
		}

		viper.AutomaticEnv() // read in environment variables that match

		// If a config file is found, read it in.
		if err := viper.ReadInConfig(); err == nil {
			fmt.Println("Using config file:", viper.ConfigFileUsed())
		}
	}
}
