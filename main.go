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

	const maxTasks = 10

	grader := grader.NewGrader(schema, processors.Processor{}, maxTasks)

	testResults := grader.Grade()

	bs, _ := json.Marshal(testResults)
	fmt.Println(string(bs))
	// fmt.Println(testResults)
}
