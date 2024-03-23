package main

import "fmt"

type Walker struct {
	Name string
}

func (w *Walker) Walk() { fmt.Printf("%s is walking\n", w.Name) }

type Runner struct {
	Walker
}

func NewRunner(name string) *Runner {
	return &Runner{Walker{Name: name}}
}

func (r *Runner) Run() { fmt.Printf("%s is running\n", r.Name) }

func main() {
	runner := NewRunner("John")
	runner.Walk()
	runner.Run()
}
