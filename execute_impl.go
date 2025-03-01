package lksh

// This repo has some interesting ideas: https://github.com/bcmills/more
// This package could be improved by using some of the ideas from this repo.
// I found it while searching for a LimitedWriter implementation
// https://github.com/golang/go/issues/51115

// Would it be beneficial to pass in io.Writers and not return a buffer?
// This would allow the caller to decide how to handle the output.

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"slices"
	"strings"

	"github.com/bcmills/more/moreio"
)

var errLimitExceeded = errors.New("buffer limit exceeded")

func execute(ctx context.Context, cfg *Config, path string, args []string) (int, *bytes.Buffer, error) {
	// Handle environment variables
	currentEnv := os.Environ()
	finalEnv := make([]string, 0)

	if len(cfg.KeepEnvVars) == 0 {
		// If KeepEnvVars slice is empty, keep all environment variables
		finalEnv = append(finalEnv, currentEnv...)
	} else {
		// If KeepEnvVars slice has entries, only keep specified variables
		for _, envVar := range currentEnv {
			parts := strings.SplitN(envVar, "=", 2)
			if len(parts) == 2 {
				varName := parts[0]
				// Check if this variable should be kept
				if slices.Contains(cfg.KeepEnvVars, varName) {
					finalEnv = append(finalEnv, envVar)
				}
			}
		}
	}

	// Add variables from AddEnvVars
	for key, value := range cfg.AddEnvVars {
		finalEnv = append(finalEnv, key+"="+value)
	}

	cmd := exec.Cmd{
		Path:  path,
		Args:  args[0:],
		Stdin: os.Stdin,
		Env:   finalEnv,
	}

	// context setup
	if ctx == nil {
		ctx = context.Background()
	}

	ctx, cancel := context.WithCancel(ctx)
	defer func() {
		// Only cancel if the context is not already done
		if ctx.Err() == nil {
			cancel()
		}
	}()

	// setup command canceling
	go func() {
		<-ctx.Done() // Block until the context is done.
		// If the context is canceled or times out, kill the process.
		if ctx.Err() != nil && cmd.Process != nil {
			fmt.Println("Context canceled or timed out. Killing the process.")
			cmd.Process.Kill()
		}
	}()

	var buf *bytes.Buffer
	var limit io.Writer
	if cfg.Pipe && cfg.MaxBufferSize <= 0 {
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
	} else if cfg.MaxBufferSize > 0 {
		buf = &bytes.Buffer{}
		limit = moreio.LimitWriter(buf, cfg.MaxBufferSize, errLimitExceeded)
		cmd.Stdout = limit
		cmd.Stderr = limit
	} else if cfg.Pipe && cfg.MaxBufferSize > 0 {
		buf = &bytes.Buffer{}
		limit = moreio.LimitWriter(buf, cfg.MaxBufferSize, errLimitExceeded)
		mw := io.MultiWriter(os.Stdout, limit)
		cmd.Stdout = mw
		cmd.Stderr = mw
	}

	err := cmd.Start()
	if err != nil {
		return cfg.DefaultExit, buf, err
	}

	err = cmd.Wait()
	if err != nil {
		return cmd.ProcessState.ExitCode(), buf, err
	}

	return cmd.ProcessState.ExitCode(), buf, nil
}
