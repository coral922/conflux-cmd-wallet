package util

import (
	"fmt"
	"os/exec"
	"runtime"
)

var browserCommands = map[string]string{
	"windows": "/c start",
	"darwin":  "open",
	"linux":   "xdg-open",
}

var baseCommands = map[string]string{
	"windows": "cmd",
	"darwin":  "bash",
	"linux":   "bash",
}

func OpenBrowser(url string) error {
	base, ok := baseCommands[runtime.GOOS]
	run, ok := browserCommands[runtime.GOOS]
	if !ok {
		return fmt.Errorf("don't know how to open things on %s platform", runtime.GOOS)
	}
	cmd := exec.Command(base, run+" "+url)
	fmt.Println(url)
	return cmd.Start()
}

func SecondConfirm(question ...string) bool {
	t := "Are you sure ?"
	if len(question) > 0 {
		t = question[0]
	}
	tt := t + " [y/n]: "
	fmt.Print(tt)
	var ans string
	_, _ = fmt.Scanln(&ans)
	if ans == "Y" || ans == "y" {
		return true
	}
	return false
}
