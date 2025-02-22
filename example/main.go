package main

import (
	"context"
	"fmt"
	"time"

	"github.com/obiwahn/go-lksh"
)

func main() {
	cfg := lksh.NewConfig()
	cfg.Pipe = true
	cfg.MaxBufferSize = lksh.MegaByte * 12
	//cfg.LookUpPath = false
	cfg.AddEnvVar("foo", "bar")

	commands := [][]string{
		{"bash", "-c", "for i in {1..4}; do echo $i; echo ${i}e >&2; sleep 1; done;"},
		{"bash", "-c", "for i in {1..10}; do echo $i; echo ${i}e >&2; sleep 1; done;"},
		{"bash", "-c", "set | grep you"},
		{"that_does_not_exist"},
		{"ls", "--color", "-lisah"},
		{"bash", "-c", "echo $PATH $foo;"},
		{"bash", "-c", "echo hello;"},
		{"bash", "-c", "cat GNUmakefile;"},
	}

	for _, cmd := range commands {

		ctx := context.Background()
		ctx, stop := context.WithTimeout(ctx, 2*time.Second)
		defer stop()

		code, buf, err := lksh.Execute(ctx, cfg, cmd)
		if err != nil {
			fmt.Println(err)
		} else {
			fmt.Printf("exit code %v\n", code)
			if buf != nil {
				fmt.Printf("buf:\n---\n%v---\n", buf)
			}
		}
	}
}
