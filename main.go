package main

import (
	"flag"
	"fmt"

	"github.com/johnhany97/grader/processors"
)

func main() {
	file := flag.String("file", "", "Name of file (required)")
	folder := flag.String("folder", "", "Folder in which the file is stored (if not in same dir as main.go)")
	// markingSchemeFile := flag.String("scheme", "", "Marking scheme to follow when grading the assignment (required)")
	// outFile := flag.String("out", "", "Where to store the results of the grader (required)")
	flag.Parse()

	// if *file == "" || *markingSchemeFile == "" || *outFile == "" {
	// 	fmt.Println("Missing arguments, please make sure to provide all required arguments.")
	// 	return
	// }

	processor := processors.Processor{}
	output := processor.Execute(*file, *folder)
	fmt.Print(output)
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}
