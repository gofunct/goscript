package script

import (
	flag "github.com/spf13/pflag"
)

var initializers []func()
var EnablePrefixMatching = false
var EnableCommandSorting = true
var minNamePadding = 11
var minCommandPathPadding = 11
var minUsagePadding = 25

type FParseErrWhitelist flag.ParseErrorsWhitelist
