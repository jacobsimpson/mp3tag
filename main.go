//go:generate pigeon -o parser/grammar.go parser/grammar.peg
//go:generate goimports -w parser/grammar.go
package main

import "github.com/jacobsimpson/mp3tag/cmd"

func main() {
	cmd.Execute()
}
