//go:build windows

package launcher

import (
	"fmt"
	"os"
	"os/exec"
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
		// Windows does not support syscall.Exec for replacing the current process image, so run synchronously and exit afterwards.
		cmd := exec.Command(commandPath, commandArgs...)
		cmd.Stdin = os.Stdin
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		if err := cmd.Run(); err != nil {
			return err
		}
		os.Exit(0)
		return nil
	}

	if path, err := exec.LookPath("wt"); err == nil {
		args := append([]string{"new-tab", commandPath}, commandArgs...)
		return launchDetachedCommand(path, args...)
	}

	args := append([]string{"/c", "start", ""}, append([]string{commandPath}, commandArgs...)...)
	return launchDetachedCommand("cmd", args...)
}

func launchDetached(v model.Variant) error {
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
	cmd.SysProcAttr = &syscall.SysProcAttr{CreationFlags: syscall.CREATE_NEW_PROCESS_GROUP}
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
