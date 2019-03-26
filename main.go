package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"

	"github.com/johnhany97/grader/test"

	"github.com/johnhany97/grader/grader"
)

func main() {
	// Parse arguments to application
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

	const maxTasks = 10

	grader := grader.NewGrader(schema, maxTasks)

	// Run Grader
	testResults := grader.Grade()

	// Process test results as requested in schema
	processResults(testResults, schema.Outfile, schema.Folder)
}

func processResults(tr []test.TestResult, outfile string, folder string) {
	bs, _ := json.Marshal(tr)
	if outfile == "" {
		fmt.Println(string(bs))
	} else {
		if err := ioutil.WriteFile(folder+outfile, bs, 0644); err != nil {
			log.Fatal(err)
		}
	}
}
