package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"strings"

	"github.com/johnhany97/grader/processors"
)

func main() {
	schemaFile := flag.String("schema", "", "Marking scheme to follow when grading the assignment (required)")
	flag.Parse()

	if *schemaFile == "" {
		fmt.Println("Missing arguments, please make sure to provide all required arguments.")
		return
	}

	var schema struct {
		File      string `json:"file"`
		Language  string `json:"language"`
		ClassName string `json:"className"`
		Folder    string `json:"folder"`
		Tests     []struct {
			Type           string   `json:"type"`
			Input          []string `json:"input"`
			ExpectedOutput []string `json:"expectedOutput"`
			UnitTest       string   `json:"unitTest"`
		}
		Outfile string `json:"outfile"`
	}
	schemaFileContent, err := ioutil.ReadFile(*schemaFile)
	check(err)
	err = json.Unmarshal(schemaFileContent, &schema)
	check(err)
	processor := processors.Processor{}

	for _, test := range schema.Tests {
		var result string
		// Obtain output
		switch test.Type {
		case "io":
			output := processor.ExecuteWithInput(schema.File, schema.Folder, strings.Join(test.Input, "\n"))
			result = fmt.Sprintf("Expected:\n%v\nGot:\n%v\nTest passed: %v", strings.Join(test.ExpectedOutput, "\n"), strings.TrimSpace(output), strings.TrimSpace(output) == strings.Join(test.ExpectedOutput, "\n"))
		case "junit":
			// Obtain Junit file shell
			junitShell, err := ioutil.ReadFile("assets/JUnit.java")
			check(err)
			// Obtain all unit tests
			// Put it all together
			junitFinal := fmt.Sprintf(string(junitShell), schema.ClassName, test.UnitTest)
			fmt.Print(processor.ExecuteJUnitTests(schema.ClassName, schema.Folder, junitFinal))
		}
		// Store output of test
		switch schema.Outfile {
		case "":
			fmt.Println(result)
		default:
			err := ioutil.WriteFile(schema.Folder+schema.Outfile, []byte(result), 0644)
			check(err)
		}
	}
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}
