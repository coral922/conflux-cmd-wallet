package main

import "cfxWorld/cmd"

//go:generate go run ./cmd/doc.go
func main() {
	cmd.Execute()
}
