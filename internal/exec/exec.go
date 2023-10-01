package exec

import (
	"context"
	"os"
	"os/exec"
	"time"
)

func RunCommand(ctx context.Context, timeout time.Duration, directory string, command string, args ...string) (string, error) {
	cmd := exec.Command(command, args...)
	cmd.Dir = directory
	//cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if timeout > 0 {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, timeout)
		defer cancel()
	}
	out, err := cmd.Output()
	if err != nil {
		return "", err
	}
	return string(out), nil
}
