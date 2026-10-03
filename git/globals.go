package git

import "regexp"

var (
	// Debug will cause all git calls to be printed out.
	Debug bool
	// SSHKey is an optional path to an SSH private key for Git
	// authentication. Useful if not using an ssh-agent.
	SSHKey string

	reRepo *regexp.Regexp = regexp.MustCompile(
		`(.+://)?([^:]+:.+@)?([^@:/]+)(:\d+)?[:/]+(.+)`,
	)
)
