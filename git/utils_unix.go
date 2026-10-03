//go:build !windows

package git

func sshCmd() string {
	return "ssh"
}
