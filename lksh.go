package lksh

import (
	"bytes"
	"context"
	"fmt"
	"os/exec"
	"strings"
	"time"

	"github.com/fatih/color"
)

func Execute(cfg *Config, args []string) (int, *bytes.Buffer, error) {
	// context setup
	ctx := context.Background()
	if cfg.Ctx != nil {
		ctx = cfg.Ctx
	}
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	// binary look up
	path := args[0]
	var err error

	if cfg.LookUpPath {
		path, err = exec.LookPath(args[0])
		if err != nil {
			return DefaultExit, nil, err
		}
	}

	// decorated
	if cfg.Decorate {
		return executeDecorated(ctx, cfg, path, args)
	} else {
		return execute(ctx, cfg, path, args)
	}
}

func executeDecorated(ctx context.Context, cfg *Config, path string, args []string) (int, *bytes.Buffer, error) {
	color.Set(color.Bold)
	color.Set(color.FgBlue)
	fmt.Printf("\n\nExecuting: %s \n\n", strings.Join(args, " "))
	color.Unset()
	defer fmt.Print("\n\n\n")

	start := time.Now()
	code, buf, err := execute(ctx, cfg, path, args)
	duration := time.Since(start)

	color.Set(color.Bold)
	if code == 0 {
		color.Set(color.FgGreen)
	} else {
		color.Set(color.FgRed)
	}
	defer color.Unset()

	fmt.Printf("\nexit with %v after %v", code, duration)
	return code, buf, err
}
