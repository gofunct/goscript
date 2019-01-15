package script

import (
	"github.com/gofunct/goscript/render/bash"
	"github.com/inconshreveable/mousetrap"
	"os"
	"time"
)

var preExecHookFn = preExecHook

func preExecHook(c *Command) {
	if bash.MousetrapHelpText != "" && mousetrap.StartedByExplorer() {
		c.Print(bash.MousetrapHelpText)
		time.Sleep(5 * time.Second)
		os.Exit(1)
	}
}
