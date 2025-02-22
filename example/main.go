package main

import (
	"context"
	"fmt"

	"github.com/obiwahn/go-lksh"
)

func main() {
	cfg := lksh.DefaultConfig()
	cfg.Ctx = context.Background()
	//cfg.LookUpPath = false
	cmd := []string{}
	cmd = []string{"./scripts/test-echo"}
	cmd = []string{"bash", "-c", "set | grep you"}
	cmd = []string{"ls", "--color", "-lisah"}
	cmd = []string{"bash", "-c", "for i in {1..10}; do echo $i; sleep 1; done;"}

	code, buf, err := lksh.Execute(cfg, cmd)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Printf("exit code %v\n", code)
		if buf != nil {
			fmt.Printf("buf:\n---\n%v---\n", buf)
		}
	}
}
