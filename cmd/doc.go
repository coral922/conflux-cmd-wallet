//+build ignore

package main

import (
	"cfxWorld/cmd"
	"log"
	"os"

	"github.com/spf13/cobra/doc"
)

const DocDir = "./doc"

func main() {
	err := os.MkdirAll(DocDir, 0666)
	if err != nil {
		log.Fatal(err)
	}
	err = doc.GenMarkdownTree(cmd.RootCmd, DocDir)
	if err != nil {
		log.Fatal(err)
	}
}
