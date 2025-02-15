package main

import (
	"gitlab.com/monokuro/era/cmd"
)

//go:generate go run gen.go

func main() {
	cmd.Execute()
}
