package src

import (
	"context"
	"fmt"
	"github.com/nsf/termbox-go"
	"io"
)

type Terminal struct {
	// The current working directory
	CurrentDir string

	// The current user
	CurrentUser string

	// The current host
	CurrentHost string

	// The current prompt
	Prompt string

	// The current context
	Context context.Context

	// History queue
	HistoryQueue *HistoryQueue
}

type KeyStroke struct {
	// The key that was pressed
	Key string
}

type Command struct {
	// The command that was entered
	Command string

	// The arguments that were passed to the command
	Args []string
}

func NewTerminal(c io.ReadWriter, prompt string) *Terminal {
	// Initilize history queue
	historyQueue := &HistoryQueue{
		Queue:       make([]string, 0),
		CurrenIndex: 0,
	}

	return &Terminal{
		Context:      context.Background(),
		Prompt:       prompt,
		HistoryQueue: historyQueue,
	}
}
func (t *Terminal) ReadLine() (cmd *Command, error error) {
	err := termbox.Init()
	cmd = &Command{}

	if err != nil {
		panic(err)
	}

	defer termbox.Close()

keyPressListerLoop:
	for {
		switch ev := termbox.PollEvent(); ev.Type {
		case termbox.EventKey:
			switch ev.Key {
			case termbox.KeyEnter:
				return cmd, nil
			case termbox.KeySpace:
				cmd.Command += " "
			case termbox.KeyBackspace, termbox.KeyBackspace2:
				cmd.Command = cmd.Command[:len(cmd.Command)-1]
			case termbox.KeyCtrlC:
				break keyPressListerLoop
			case termbox.KeyCtrlD:
				break keyPressListerLoop
			case termbox.KeyArrowUp:
				// TODO: Implement history
				fmt.Println("Up")
			case termbox.KeyArrowDown:
				// TODO: Implement history
				fmt.Println("Down")

			default:
				cmd.Command += string(ev.Ch)
			}

			termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)
			termbox.SetCursor(len(cmd.Command), 0)
			termbox.Flush()

		}
	}

	return cmd, nil
}
