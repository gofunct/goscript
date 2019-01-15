package utils

import (
	"github.com/spf13/pflag"
)

func HasNoOptDefVal(name string, fs *pflag.FlagSet) bool {
	flag := fs.Lookup(name)
	if flag == nil {
		return false
	}
	return flag.NoOptDefVal != ""
}

func ShortHasNoOptDefVal(name string, fs *pflag.FlagSet) bool {
	if len(name) == 0 {
		return false
	}

	flag := fs.ShorthandLookup(name[:1])
	if flag == nil {
		return false
	}
	return flag.NoOptDefVal != ""
}
