package main

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

func main() {
	const (
		prompt = "\033[01;32m[%s@%s\033[0m \033[01;37m%s\033[01;32m] $\033[0m "
	)

	var (
		inputArr      []string
		inputCmd      string
		parentContext context.Context
		cmdContext    context.Context
		cmdCurDirr    string
		currUser      string
		currHost      string
	)
	_ = cmdContext
	_ = inputCmd

	for {
		cmdCurDirr = os.Getenv("PWD")
		currUser = os.Getenv("USER")
		currHost = os.Getenv("HOSTNAME")

		fmt.Printf(prompt, currUser, currHost, cmdCurDirr)

		reader := bufio.NewReader(os.Stdin)
		inputCmd, err := reader.ReadString('\n')

		if err != nil {
			fmt.Println(err)
			continue
		}

		inputArr = strings.Fields(inputCmd)

		// create a context that
		parentContext = context.Background()
		cmdContext, cancelCmdCtx := context.WithCancel(parentContext)

		// Create a new subprocess
		subProcess := exec.CommandContext(cmdContext, inputArr[0], inputArr[1:]...)

		// Set dir
		subProcess.Dir = cmdCurDirr

		// Create
		stdIn, err := subProcess.StdinPipe()

		// check runtime error
		if err != nil {
			fmt.Println(err)
			continue
		}

		subProcess.Stdin = os.Stdin
		subProcess.Stdout = os.Stdout
		subProcess.Stderr = os.Stderr

		// Start the process
		if err = subProcess.Start(); err != nil {
			continue
		}

		err = subProcess.Wait()
		if err != nil {
			fmt.Println(err)
			continue
		}
		cancelCmdCtx()
		stdIn.Close()
	}

}
