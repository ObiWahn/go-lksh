package lksh

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"os"
	"os/exec"
)

func execute(ctx context.Context, cfg *Config, path string, args []string) (int, *bytes.Buffer, error) {
	cmd := exec.Cmd{
		Path: path,
		Args: args[0:],

		Stdin:  os.Stdin,
		Stderr: os.Stderr,
	}

	// setup command canceling
	go func() {
		<-ctx.Done() // Block until the context is done.
		// If the context is canceled or times out, kill the process.
		if ctx.Err() != nil && cmd.Process != nil {
			fmt.Println("Context canceled or timed out. Killing the process.")
			cmd.Process.Kill()
		}
	}()

	var stdoutPipe io.ReadCloser
	var err error
	var buf *bytes.Buffer

	if cfg.MaxBufferSize > 0 {
		buf = &bytes.Buffer{}
		stdoutPipe, err = cmd.StdoutPipe()
		if err != nil {
			return DefaultExit, buf, err
		}
	} else if cfg.Pipe {
		cmd.Stdout = os.Stdout
	}

	err = cmd.Start()
	if err != nil {
		return DefaultExit, buf, err
	}

	if cfg.MaxBufferSize > 0 {
		go func() {
			// Limit reader to 2GB to prevent excessive memory usage
			limitedReader := &io.LimitedReader{R: stdoutPipe, N: cfg.MaxBufferSize}
			if cfg.Pipe {
				// Create a MultiWriter (writes to buffer + stdout)
				multiWriter := io.MultiWriter(buf, os.Stdout) // Comment out os.Stdout if not needed
				io.Copy(multiWriter, limitedReader)
			} else {
				io.Copy(buf, limitedReader)
			}
		}()
	}

	err = cmd.Wait()
	if err != nil {
		return cmd.ProcessState.ExitCode(), buf, err
	}

	return cmd.ProcessState.ExitCode(), buf, nil
}
