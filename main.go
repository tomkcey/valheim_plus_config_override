package main

import (
	"fmt"
	"lib/input"
	"lib/processor"
	"os"
)

func main() {
	store := input.Prepare()
	if store == nil {
		os.Exit(0)
	}
	sm := *processor.PreProcess(store.Source)
	if sm == nil {
		fmt.Println("Given source file was somehow detected but nothing came of it.")
		os.Exit(20)
	}
	tm := *processor.PreProcess(store.Target)
	if tm == nil {
		fmt.Println("Given target file was somehow detected but nothing came of it.")
		os.Exit(20)
	}

	m := processor.Process(sm, tm)

	processor.PostProcess(m)
}
