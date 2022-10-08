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
		multiCommands [][]string
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
		var err error
		cmdCurDirr, err = os.Getwd()
		if err != nil {
			panic(err)
		}
		currUser = os.Getenv("USER")
		currHost = os.Getenv("HOSTNAME")

		fmt.Printf(prompt, currUser, currHost, cmdCurDirr)
		reader := bufio.NewReader(os.Stdin)
		inputCmd, err := reader.ReadString('\n')

		if err != nil {
			fmt.Println(err)
			continue
		}

		if inputCmd == "exit\n" {
			fmt.Println("Bye!")
			os.Exit(0)
		}

		// if there is no input, continue
		if inputCmd == "\n" {
			continue
		}

		multiCommands = make([][]string, 0)

		// detect if there is a multi command
		if strings.Contains(inputCmd, ";") {
			inputArr = strings.Split(inputCmd, ";")
			for _, cmd := range inputArr {
				multiCommands = append(multiCommands, strings.Fields(cmd))
			}
		} else {
			multiCommands = append(multiCommands, strings.Fields(inputCmd))
		}

		for _, v := range multiCommands {
			// detect variable assignment
			var (
				subCmd []string = v
			)

			if strings.Contains(subCmd[0], "=") {
				variable := strings.Split(v[0], "=")
				os.Setenv(variable[0], variable[1])
				continue
			}

			// detect variables and replace them
			for i, v := range v {
				if strings.HasPrefix(v, "$") {
					subCmd[i] = os.Getenv(v[1:])
				}
			}

			// if there is no input, continue
			if len(v) == 0 {
				continue
			}

			if v[0] == "cd" {
				if len(subCmd) == 1 {
					os.Chdir(os.Getenv("HOME"))
					continue
				}

				err := os.Chdir(subCmd[1])
				if err != nil {
					fmt.Println(err)
					continue
				}
			}

			// detect if there is a pipe
			if strings.Contains(inputCmd, "|") {

			}

			// create a context that
			parentContext = context.Background()
			cmdContext, cancelCmdCtx := context.WithCancel(parentContext)

			// Create a new subprocess
			subProcess := exec.CommandContext(cmdContext, subCmd[0], subCmd[1:]...)

			// Set dir
			subProcess.Dir = cmdCurDirr

			// Create
			stdIn, err := subProcess.StdinPipe()

			// check runtime error
			if err != nil {
				fmt.Fprintln(os.Stderr, err)
				continue
			}

			subProcess.Stdin = os.Stdin
			subProcess.Stdout = os.Stdout
			subProcess.Stderr = os.Stderr

			// Start the process
			if err = subProcess.Start(); err != nil {
				errStr := err.Error()
				errStr = strings.Replace(errStr, "exec: ", "osh: ", 1)
				// write to os.Stderr
				fmt.Fprintln(os.Stderr, errStr)
				continue
			}

			err = subProcess.Wait()
			if err != nil {
				// write to os.Stderr
				fmt.Fprintln(os.Stderr, err)
				continue
			}
			cancelCmdCtx()
			stdIn.Close()
		}
	}

}
