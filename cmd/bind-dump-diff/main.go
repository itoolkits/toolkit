// bind dump db diff
// ./rndc dumpdb -zones
// print diff table, use go-pretty

package main

import (
	"fmt"
	"os"

	"github.com/itoolkits/toolkit/dnt"
)

// process main function
func main() {
	args := os.Args[1:]

	if len(args) != 2 {
		fmt.Printf("only support 2 args")
		os.Exit(1)
	}

	A, err := dnt.ParseDumpDB(args[0])
	if err != nil {
		fmt.Printf("parse dumpdb error,%s %v", args[0], err)
		os.Exit(1)
	}
	B, err := dnt.ParseDumpDB(args[1])
	if err != nil {
		fmt.Printf("parse dumpdb error,%s %v", args[1], err)
		os.Exit(1)
	}

	handler := NewDiffHandler(A, B)
	handler.Start()
	handler.PrintDiff(args[0], args[1])
}
