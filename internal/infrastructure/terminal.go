package infrastructure

import (
	"os/exec"
	"start/internal/config"
	"strings"
)

type Terminal struct {
	UseSudo bool
}

func NewTerminal(cfg *config.TerminalConfig) *Terminal {
	terminal := Terminal{
		UseSudo: cfg.UseSudo,
	}

	return &terminal
}

func (t *Terminal) ExecuteCmd(cmd string, dir string) (string, error) {
	if t.UseSudo {
		cmd = "sudo " + cmd
	}

	parts := strings.Fields(cmd)
	cmdExec := exec.Command(parts[0], parts[1:]...)
	cmdExec.Dir = dir

	out, err := cmdExec.CombinedOutput()
	return string(out), err
}
