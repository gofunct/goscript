package script

// InitDefaultVersionFlag adds default version flag to c.
// It is called automatically by executing the c.
// If c already has a version flag, it will do nothing.
// If c.Version is empty, it will do nothing.
func (c *Command) InitDefaultVersionFlag() {
	if c.Config.Version == "" {
		return
	}

	c.mergePersistentFlags()
	if c.Flags().Lookup("version") == nil {
		usage := "version for "
		if c.Name() == "" {
			usage += "this command"
		} else {
			usage += c.Name()
		}
		c.Flags().Bool("version", false, usage)
	}
}

// SetVersionTemplate sets version template to be used. Application can use it to set custom template.
func (c *Command) SetVersionTemplate(s string) {
	c.versionTemplate = s
}

// VersionTemplate return version template for the command.
func (c *Command) VersionTemplate() string {
	if c.versionTemplate != "" {
		return c.versionTemplate
	}

	if c.HasParent() {
		return c.parent.VersionTemplate()
	}
	return `{{with .Name}}{{printf "%s " .}}{{end}}{{printf "version %s" .Version}}
`
}
