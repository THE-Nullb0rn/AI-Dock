//go:build linux

package launcher

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"syscall"

	"github.com/THE-Nullb0rn/AI-Dock/internal/model"
)

func launchCLI(v model.Variant, mode string) error {
	commandPath, commandArgs, err := commandForVariant(v)
	if err != nil {
		return err
	}

	if mode == "same_terminal" {
		// syscall.Exec replaces the current process image, which is the closest native behavior for reuse of the same terminal.
		return syscall.Exec(commandPath, append([]string{commandPath}, commandArgs...), os.Environ())
	}

	terminalPath, terminalArgs, err := terminalCommand(commandPath, commandArgs)
	if err != nil {
		return err
	}
	return launchDetachedCommand(terminalPath, terminalArgs...)
}

func launchDetached(v model.Variant) error {
	if strings.HasSuffix(strings.ToLower(v.Path), ".desktop") {
		if path, err := exec.LookPath("gtk-launch"); err == nil {
			name := strings.TrimSuffix(filepath.Base(v.Path), filepath.Ext(v.Path))
			return launchDetachedCommand(path, name)
		}
	}

	commandPath, commandArgs, err := commandForVariant(v)
	if err != nil {
		return err
	}
	return launchDetachedCommand(commandPath, commandArgs...)
}

func launchDetachedCommand(command string, args ...string) error {
	cmd := exec.Command(command, args...)
	configureDetachedCommand(cmd)

	devNull, err := os.OpenFile(os.DevNull, os.O_RDWR, 0)
	if err != nil {
		return fmt.Errorf("open devnull: %w", err)
	}
	defer devNull.Close()

	cmd.Stdin = devNull
	cmd.Stdout = devNull
	cmd.Stderr = devNull
	return cmd.Start()
}

func configureDetachedCommand(cmd *exec.Cmd) {
	cmd.SysProcAttr = &syscall.SysProcAttr{Setsid: true}
}

func commandForVariant(v model.Variant) (string, []string, error) {
	if strings.TrimSpace(v.LaunchCmd) != "" {
		return parseCommandLine(v.LaunchCmd)
	}
	if strings.TrimSpace(v.Path) != "" {
		return parseCommandLine(v.Path)
	}
	return "", nil, fmt.Errorf("variant %q has no launch command", v.Type)
}

func terminalCommand(command string, args []string) (string, []string, error) {
	terminalName := strings.TrimSpace(os.Getenv("TERMINAL"))
	if terminalName != "" {
		if terminalPath, err := exec.LookPath(terminalName); err == nil {
			return terminalPath, append([]string{"-e", command}, args...), nil
		}
	}

	type terminalSpec struct {
		name  string
		flags []string
	}

	for _, spec := range []terminalSpec{
		{name: "kitty", flags: []string{"-e"}},
		{name: "alacritty", flags: []string{"-e"}},
		{name: "gnome-terminal", flags: []string{"--"}},
		{name: "konsole", flags: []string{"-e"}},
		{name: "xterm", flags: []string{"-e"}},
	} {
		if terminalPath, err := exec.LookPath(spec.name); err == nil {
			return terminalPath, append(spec.flags, append([]string{command}, args...)...), nil
		}
	}

	return "", nil, fmt.Errorf("no supported terminal emulator found")
}
