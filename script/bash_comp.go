package script

import (
	"bytes"
	"fmt"
	"github.com/gofunct/goscript/render/bash"
	"github.com/spf13/pflag"
	"os"
	"sort"
	"strings"
	"io"
)

func writeCommands(buf *bytes.Buffer, cmd *Command) {
	buf.WriteString("    commands=()\n")
	for _, c := range cmd.Commands() {
		if !c.IsAvailableCommand() || c == cmd.helpCommand {
			continue
		}
		buf.WriteString(fmt.Sprintf("    commands+=(%q)\n", c.Name()))
		writeCmdAliases(buf, c)
	}
	buf.WriteString("\n")
}

func writeFlagHandler(buf *bytes.Buffer, name string, annotations map[string][]string, cmd *Command) {
	for key, value := range annotations {
		switch key {
		case bash.BashCompFilenameExt:
			buf.WriteString(fmt.Sprintf("    flags_with_completion+=(%q)\n", name))

			var ext string
			if len(value) > 0 {
				ext = fmt.Sprintf("__%s_handle_filename_extension_flag ", cmd.Root().Name()) + strings.Join(value, "|")
			} else {
				ext = "_filedir"
			}
			buf.WriteString(fmt.Sprintf("    flags_completion+=(%q)\n", ext))
		case bash.BashCompCustom:
			buf.WriteString(fmt.Sprintf("    flags_with_completion+=(%q)\n", name))
			if len(value) > 0 {
				handlers := strings.Join(value, "; ")
				buf.WriteString(fmt.Sprintf("    flags_completion+=(%q)\n", handlers))
			} else {
				buf.WriteString("    flags_completion+=(:)\n")
			}
		case bash.BashCompSubdirsInDir:
			buf.WriteString(fmt.Sprintf("    flags_with_completion+=(%q)\n", name))

			var ext string
			if len(value) == 1 {
				ext = fmt.Sprintf("__%s_handle_subdirs_in_dir_flag ", cmd.Root().Name()) + value[0]
			} else {
				ext = "_filedir -d"
			}
			buf.WriteString(fmt.Sprintf("    flags_completion+=(%q)\n", ext))
		}
	}
}

func writeShortFlag(buf *bytes.Buffer, flag *pflag.Flag, cmd *Command) {
	name := flag.Shorthand
	format := "    "
	if len(flag.NoOptDefVal) == 0 {
		format += "two_word_"
	}
	format += "flags+=(\"-%s\")\n"
	buf.WriteString(fmt.Sprintf(format, name))
	writeFlagHandler(buf, "-"+name, flag.Annotations, cmd)
}

func writeFlag(buf *bytes.Buffer, flag *pflag.Flag, cmd *Command) {
	name := flag.Name
	format := "    flags+=(\"--%s"
	if len(flag.NoOptDefVal) == 0 {
		format += "="
	}
	format += "\")\n"
	buf.WriteString(fmt.Sprintf(format, name))
	writeFlagHandler(buf, "--"+name, flag.Annotations, cmd)
}

func writeLocalNonPersistentFlag(buf *bytes.Buffer, flag *pflag.Flag) {
	name := flag.Name
	format := "    local_nonpersistent_flags+=(\"--%s"
	if len(flag.NoOptDefVal) == 0 {
		format += "="
	}
	format += "\")\n"
	buf.WriteString(fmt.Sprintf(format, name))
}

func writeFlags(buf *bytes.Buffer, cmd *Command) {
	buf.WriteString(`    flags=()
    two_word_flags=()
    local_nonpersistent_flags=()
    flags_with_completion=()
    flags_completion=()
`)
	localNonPersistentFlags := cmd.LocalNonPersistentFlags()
	cmd.NonInheritedFlags().VisitAll(func(flag *pflag.Flag) {
		if nonCompletableFlag(flag) {
			return
		}
		writeFlag(buf, flag, cmd)
		if len(flag.Shorthand) > 0 {
			writeShortFlag(buf, flag, cmd)
		}
		if localNonPersistentFlags.Lookup(flag.Name) != nil {
			writeLocalNonPersistentFlag(buf, flag)
		}
	})
	cmd.InheritedFlags().VisitAll(func(flag *pflag.Flag) {
		if nonCompletableFlag(flag) {
			return
		}
		writeFlag(buf, flag, cmd)
		if len(flag.Shorthand) > 0 {
			writeShortFlag(buf, flag, cmd)
		}
	})

	buf.WriteString("\n")
}

func writeRequiredFlag(buf *bytes.Buffer, cmd *Command) {
	buf.WriteString("    must_have_one_flag=()\n")
	flags := cmd.NonInheritedFlags()
	flags.VisitAll(func(flag *pflag.Flag) {
		if nonCompletableFlag(flag) {
			return
		}
		for key := range flag.Annotations {
			switch key {
			case bash.BashCompOneRequiredFlag:
				format := "    must_have_one_flag+=(\"--%s"
				if flag.Value.Type() != "bool" {
					format += "="
				}
				format += "\")\n"
				buf.WriteString(fmt.Sprintf(format, flag.Name))

				if len(flag.Shorthand) > 0 {
					buf.WriteString(fmt.Sprintf("    must_have_one_flag+=(\"-%s\")\n", flag.Shorthand))
				}
			}
		}
	})
}

func writeRequiredNouns(buf *bytes.Buffer, cmd *Command) {
	buf.WriteString("    must_have_one_noun=()\n")
	sort.Sort(sort.StringSlice(cmd.ValidArgs))
	for _, value := range cmd.ValidArgs {
		buf.WriteString(fmt.Sprintf("    must_have_one_noun+=(%q)\n", value))
	}
}

func writeCmdAliases(buf *bytes.Buffer, cmd *Command) {
	if len(cmd.Aliases) == 0 {
		return
	}

	sort.Sort(sort.StringSlice(cmd.Aliases))

	buf.WriteString(fmt.Sprint(`    if [[ -z "${BASH_VERSION}" || "${BASH_VERSINFO[0]}" -gt 3 ]]; then`, "\n"))
	for _, value := range cmd.Aliases {
		buf.WriteString(fmt.Sprintf("        command_aliases+=(%q)\n", value))
		buf.WriteString(fmt.Sprintf("        aliashash[%q]=%q\n", value, cmd.Name()))
	}
	buf.WriteString(`    fi`)
	buf.WriteString("\n")
}
func writeArgAliases(buf *bytes.Buffer, cmd *Command) {
	buf.WriteString("    noun_aliases=()\n")
	sort.Sort(sort.StringSlice(cmd.ArgAliases))
	for _, value := range cmd.ArgAliases {
		buf.WriteString(fmt.Sprintf("    noun_aliases+=(%q)\n", value))
	}
}

func gen(buf *bytes.Buffer, cmd *Command) {
	for _, c := range cmd.Commands() {
		if !c.IsAvailableCommand() || c == cmd.helpCommand {
			continue
		}
		gen(buf, c)
	}
	commandName := cmd.CommandPath()
	commandName = strings.Replace(commandName, " ", "_", -1)
	commandName = strings.Replace(commandName, ":", "__", -1)

	if cmd.Root() == cmd {
		buf.WriteString(fmt.Sprintf("_%s_root_command()\n{\n", commandName))
	} else {
		buf.WriteString(fmt.Sprintf("_%s()\n{\n", commandName))
	}

	buf.WriteString(fmt.Sprintf("    last_command=%q\n", commandName))
	buf.WriteString("\n")
	buf.WriteString("    command_aliases=()\n")
	buf.WriteString("\n")

	writeCommands(buf, cmd)
	writeFlags(buf, cmd)
	writeRequiredFlag(buf, cmd)
	writeRequiredNouns(buf, cmd)
	writeArgAliases(buf, cmd)
	buf.WriteString("}\n\n")
}

// GenBashCompletion generates bash completion file and writes to the passed writer.
func (c *Command) GenBashCompletion(w io.Writer) error {
	buf := new(bytes.Buffer)
	bash.WritePreamble(buf, c.Name())
	if len(c.BashCompletionFunction) > 0 {
		buf.WriteString(c.BashCompletionFunction + "\n")
	}
	gen(buf, c)
	bash.WritePostscript(buf, c.Name())

	_, err := buf.WriteTo(w)
	return err
}

func nonCompletableFlag(flag *pflag.Flag) bool {
	return flag.Hidden || len(flag.Deprecated) > 0
}

// GenBashCompletionFile generates bash completion file.
func (c *Command) GenBashCompletionFile(filename string) error {
	outFile, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer outFile.Close()

	return c.GenBashCompletion(outFile)
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
	return flags.SetAnnotation(name, bash.BashCompOneRequiredFlag, []string{"true"})
}

// MarkFlagFilename adds the BashCompFilenameExt annotation to the named flag, if it exists.
// Generated bash autocompletion will select filenames for the flag, limiting to named extensions if provided.
func (c *Command) MarkFlagFilename(name string, extensions ...string) error {
	return MarkFlagFilename(c.Flags(), name, extensions...)
}

// MarkFlagCustom adds the BashCompCustom annotation to the named flag, if it exists.
// Generated bash autocompletion will call the bash function f for the flag.
func (c *Command) MarkFlagCustom(name string, f string) error {
	return MarkFlagCustom(c.Flags(), name, f)
}

// MarkPersistentFlagFilename adds the BashCompFilenameExt annotation to the named persistent flag, if it exists.
// Generated bash autocompletion will select filenames for the flag, limiting to named extensions if provided.
func (c *Command) MarkPersistentFlagFilename(name string, extensions ...string) error {
	return MarkFlagFilename(c.PersistentFlags(), name, extensions...)
}

// MarkFlagFilename adds the BashCompFilenameExt annotation to the named flag in the flag set, if it exists.
// Generated bash autocompletion will select filenames for the flag, limiting to named extensions if provided.
func MarkFlagFilename(flags *pflag.FlagSet, name string, extensions ...string) error {
	return flags.SetAnnotation(name, bash.BashCompFilenameExt, extensions)
}

// MarkFlagCustom adds the BashCompCustom annotation to the named flag in the flag set, if it exists.
// Generated bash autocompletion will call the bash function f for the flag.
func MarkFlagCustom(flags *pflag.FlagSet, name string, f string) error {
	return flags.SetAnnotation(name, bash.BashCompCustom, []string{f})
}
