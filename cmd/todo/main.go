package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/mxygem/todo"
)

const todoFileName = ".todo.json"

func main() {
	flag.Usage = func() {
		fmt.Fprintf(flag.CommandLine.Output(),
			"%s tool.\ndeveloped while working through the book \"Powerful Command-Line Applications in Go\"\n",
			os.Args[0])
		fmt.Fprintln(flag.CommandLine.Output(), "usage:")
		flag.PrintDefaults()
	}

	task := flag.String("task", "", "task to be included in the todo list")
	list := flag.Bool("list", false, "list all tasks")
	complete := flag.Int("complete", 0, "item to be completed")
	flag.Parse()

	l := &todo.List{}
	if err := l.Get(todoFileName); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	switch {
	case *list:
		fmt.Print(l)
	case *complete > 0:
		if err := l.Complete(*complete); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}

		if err := l.Save(todoFileName); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
	case *task != "":
		l.Add(*task)

		if err := l.Save(todoFileName); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
	default:
		fmt.Fprintln(os.Stderr, "invalid option")
		os.Exit(1)
	}
}
