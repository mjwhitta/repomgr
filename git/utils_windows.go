//go:build windows

package git

import (
	"bytes"
	"os/exec"
)

func sshCmd() string {
	var cmd *exec.Cmd = exec.Command(
		"git",
		"config",
		"core.sshCommand",
	)

	if b, e := cmd.CombinedOutput(); e == nil {
		if b = bytes.TrimSpace(b); len(b) > 0 {
			return string(b)
		}
	}

	return "C:/Windows/System32/OpenSSH/ssh.exe"
}
