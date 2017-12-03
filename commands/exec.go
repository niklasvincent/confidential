package commands

import (
	"fmt"
	"os"
	"os/exec"
	"syscall"
	"os/signal"

	"github.com/urfave/cli"
	"github.com/nlindblad/confidential/environment"
	"github.com/pkg/errors"
)

type commandWithFlags struct {
	executable string
	arguments []string
}

func extractCommandWithFlags(rawArgs []string) (*commandWithFlags, error) {
	start := -1
	for i, s := range rawArgs {
		if s == "--" {
			start = i + 1
			break
		}
	}

	if start == -1  || start + 1 > len(rawArgs) {
		return nil, fmt.Errorf("no command provided")
	}

	return &commandWithFlags{executable: rawArgs[start], arguments: rawArgs[start + 1:]}, nil
}

func getEnvironmentFromContext(c *cli.Context) (*environment.Environment, error) {
	retrievedEnvironment, err := RetrieveEnvironmentVariablesFromContext(c)
	if err != nil {
		return nil, err
	}

	runtimeEnvironment, err := environment.NewEnvironmentFromRuntime(os.Environ())
	if err != nil {
		return nil, err
	}

	newEnvironment := runtimeEnvironment.Union(retrievedEnvironment)

	return newEnvironment, nil
}

func execRun(c *cli.Context) error {
	curatedEnvironment, err := getEnvironmentFromContext(c)
	if err != nil {
		return err
	}

	args := os.Args[1:]
	cmd, err := extractCommandWithFlags(args)
	if err != nil {
		return err
	}

	newCmd := exec.Command(cmd.executable, cmd.arguments...)
	newCmd.Stdin = os.Stdin
	newCmd.Stdout = os.Stdout
	newCmd.Stderr = os.Stderr
	newCmd.Env = curatedEnvironment.AsStrings()

	// Forward SIGINT, SIGTERM, SIGKILL to the child command
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGTERM, os.Interrupt, os.Kill)

	go func() {
		sig := <-sigChan
		if newCmd.Process != nil {
			newCmd.Process.Signal(sig)
		}
	}()

	var waitStatus syscall.WaitStatus
	if err := newCmd.Run(); err != nil {
		if err != nil {
			return errors.Wrap(err, "failed to run command")
		}
		if exitError, ok := err.(*exec.ExitError); ok {
			waitStatus = exitError.Sys().(syscall.WaitStatus)
			os.Exit(waitStatus.ExitStatus())
		}
	}
	return nil
}