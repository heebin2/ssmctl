package ssm

import (
	"fmt"
	"os"
	"os/exec"
)

func startSession(target, user string) error {
	cmd := exec.Command(
		"aws",
		"ssm",
		"start-session",
		"--target", target,
		"--document-name", "AWS-StartInteractiveCommand",
		"--parameters", fmt.Sprintf("command=sudo su - %s", user),
	)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("start session: %w", err)
	}
	return nil
}
