package main

import (
	"fmt"
	"gorilla/repl"
	"os"
	"runtime"
)

func main() {
	// user, err := user.Current()
	// if err != nil {
	// 	panic(err)
	// }

	fmt.Printf(
		"Gorilla Programming Lanugage version 0.1 [%s]\n",
		runtime.GOOS,
	)
	repl.Start(os.Stdin, os.Stdout)
}
