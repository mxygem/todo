package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/mxygem/todo"
)

var todoFileName = ".todo.json"

func main() {
	flag.Usage = func() {
		fmt.Fprintf(flag.CommandLine.Output(),
			"%s tool.\ndeveloped while working through the book \"Powerful Command-Line Applications in Go\"\n",
			os.Args[0])
		fmt.Fprintln(flag.CommandLine.Output(), "usage:")
		flag.PrintDefaults()
	}

	add := flag.Bool("add", false, "add a task to a todo list")
	del := flag.Int("del", 0, "delete an item from a todo list")
	list := flag.Bool("list", false, "list all tasks")
	complete := flag.Int("complete", 0, "item to be completed")
	flag.Parse()

	if os.Getenv("TODO_FILENAME") != "" {
		todoFileName = os.Getenv("TODO_FILENAME")
	}

	l := &todo.List{}
	if err := l.Get(todoFileName); err != nil {
		cliError(err)
	}

	switch {
	case *list:
		fmt.Print(l)
	case *complete > 0:
		if err := l.Complete(*complete); err != nil {
			cliError(err)
		}

		if err := l.Save(todoFileName); err != nil {
			cliError(err)
		}
	case *add:
		t, err := getTask(os.Stdin, flag.Args()...)
		if err != nil {
			cliError(err)
		}
		l.Add(t)

		if err := l.Save(todoFileName); err != nil {
			cliError(err)
		}
	case *del > 0:
		if err := l.Delete(*del); err != nil {
			cliError(err)
		}

		if err := l.Save(todoFileName); err != nil {
			cliError(err)
		}
	default:
		fmt.Fprintln(os.Stderr, "invalid option")
		os.Exit(1)
	}
}

func getTask(r io.Reader, args ...string) (string, error) {
	if len(args) > 0 {
		return strings.Join(args, " "), nil
	}

	s := bufio.NewScanner(r)
	s.Scan()
	if err := s.Err(); err != nil {
		return "", err
	}
	if len(s.Text()) == 0 {
		return "", fmt.Errorf("task cannot be blank")
	}

	return s.Text(), nil
}

func cliError(err error) {
	fmt.Fprintln(os.Stderr, err)
	os.Exit(1)
}
