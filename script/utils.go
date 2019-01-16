package script

import (
	"github.com/spf13/pflag"
)

var initializers []func()
var EnablePrefixMatching = false
var EnableCommandSorting = true
var minNamePadding = 11
var minCommandPathPadding = 11
var minUsagePadding = 25

type FParseErrWhitelist pflag.ParseErrorsWhitelist

// Sorts commands by their names.
type commandSorterByName []*Command

func (c commandSorterByName) Len() int           { return len(c) }
func (c commandSorterByName) Swap(i, j int)      { c[i], c[j] = c[j], c[i] }
func (c commandSorterByName) Less(i, j int) bool { return c[i].Name() < c[j].Name() }
