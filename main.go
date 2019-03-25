package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"

	"github.com/johnhany97/grader/grader"
	"github.com/johnhany97/grader/processors"
)

func main() {
	schemaFile := flag.String("schema", "", "Marking scheme to follow when grading the assignment (required)")
	flag.Parse()

	if *schemaFile == "" {
		fmt.Println("Missing arguments, please make sure to provide all required arguments.")
		return
	}

	schema := grader.Schema{}

	// read schema file
	schemaFileContent, err := ioutil.ReadFile(*schemaFile)
	if err != nil {
		log.Fatal(err)
	}

	// unmarshal unto the schema var
	err = json.Unmarshal(schemaFileContent, &schema)
	if err != nil {
		log.Fatal(err)
	}

	grader := grader.NewGrader(schema, processors.Processor{})

	testResults := grader.Grade()

	bs, _ := json.Marshal(testResults)
	fmt.Println(string(bs))
	// fmt.Println(testResults)
}

// func oldMain() {
// 	schemaFile := flag.String("schema", "", "Marking scheme to follow when grading the assignment (required)")
// 	flag.Parse()

// 	if *schemaFile == "" {
// 		fmt.Println("Missing arguments, please make sure to provide all required arguments.")
// 		return
// 	}

// 	var schema Schema

// 	schemaFileContent, err := ioutil.ReadFile(*schemaFile)
// 	check(err)
// 	err = json.Unmarshal(schemaFileContent, &schema)
// 	check(err)

// 	testResults := []testResult{}

// 	processor := processors.Processor{}

// 	for _, test := range schema.Tests {
// 		var result string
// 		// Obtain output
// 		switch test.Type {
// 		case "io":
// 			output, errorOutput := processor.ExecuteWithInput(schema.File, schema.Folder, strings.Join(test.Input, "\n"))
// 			testResults = append(testResults, NewTestResult(&test, output, errorOutput))
// 		case "junit":
// 			// Obtain Junit file shell
// 			junitShell, err := ioutil.ReadFile("assets/JUnit.java")
// 			check(err)
// 			// Obtain all unit tests
// 			// Put it all together
// 			junitFinal := fmt.Sprintf(string(junitShell), schema.ClassName, test.UnitTest)
// 			result, err = processor.ExecuteJUnitTests(schema.ClassName, schema.Folder, junitFinal)
// 			check(err)
// 			testResults = append(testResults, NewTestResult(&test, result, ""))
// 		}
// 		// Store output of test
// 		switch schema.Outfile {
// 		case "":
// 			fmt.Println(result)
// 		default:
// 			err := ioutil.WriteFile(schema.Folder+schema.Outfile, []byte(result), 0644)
// 			check(err)
// 		}
// 		fmt.Printf("%v", testResults)
// 	}
// }

// func check(e error) {
// 	if e != nil {
// 		panic(e)
// 	}
// }
