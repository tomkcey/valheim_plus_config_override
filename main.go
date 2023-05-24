package main

import (
	"fmt"
	"os"
	"time"

	input "example.com/tomkcey/m/input"
	postprocessor "example.com/tomkcey/m/postprocessor"
	preprocessor "example.com/tomkcey/m/preprocessor"
	processor "example.com/tomkcey/m/processor"
)

func main() {
	store := input.Prepare()
	if store == nil {
		fmt.Println("A problem occured.")
		os.Exit(20)
	}

	now := time.Now()

	ms := preprocessor.PreProcess(store)
	m := processor.Process(ms)
	postprocessor.PostProcess(*m)

	fmt.Println(time.Since(now))
}
