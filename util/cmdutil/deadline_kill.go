package cmdutil

import (
	"os/exec"
	"time"
)

// DeadlineKill kills the given exec.Cmd if execution takes longer than the timeout. The command must already be running
// when calling this function. Any error from the cmd (if force killed or not) will be returned.
func DeadlineKill(cmd *exec.Cmd, timeout time.Duration) error {
	done := make(chan error, 1)
	go func() {
		done <- cmd.Wait()
	}()

	select {
	case <-time.After(timeout):
		err := cmd.Process.Kill()
		if err != nil {
			return err
		}
		return <-done
	case err := <-done:
		return err
	}
}
