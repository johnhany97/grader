package grader

import (
	"fmt"
	"sync"

	"github.com/johnhany97/grader/processors"
	"github.com/johnhany97/grader/test"
)

type Schema struct {
	File      string      `json:"file"`
	Language  string      `json:"language"`
	ClassName string      `json:"className"`
	Folder    string      `json:"folder"`
	Tests     []test.Test `json:"tests"`
	Outfile   string      `json:"outfile"`
}

type Grader struct {
	Schema    Schema               `json:"schema"`
	Processor processors.Processor `json:"processor"`
}

const maxTasks = 100

var wg sync.WaitGroup

func NewGrader(s Schema, p processors.Processor) Grader {
	return Grader{
		Schema:    s,
		Processor: p,
	}
}

func (g Grader) Grade() []test.TestResult {
	// make a channel with a capacity of 100.
	jobChan := make(chan test.TestTask, maxTasks)

	// results
	testresults := []test.TestResult{}
	// start the worker
	wg.Add(1)
	go worker(jobChan, &testresults)

	for i := 0; i < len(g.Schema.Tests); i++ {
		// get test
		t := g.Schema.Tests[i]
		// create test task
		var task test.TestTask
		switch t.Type {
		case "output":
			task = test.OutputTestHandler{
				Test:   t,
				File:   g.Schema.File,
				Folder: g.Schema.Folder,
			}
		case "io":
			task = test.InputOutputTestHandler{
				Test:   t,
				File:   g.Schema.File,
				Folder: g.Schema.Folder,
			}
		case "junit":
			task = test.JUnitTestHandler{
				Test:      t,
				File:      g.Schema.File,
				Folder:    g.Schema.Folder,
				ClassName: g.Schema.ClassName,
			}
		}
		// enqueue task
		jobChan <- task
	}

	// we sent out all test tasks
	close(jobChan)

	wg.Wait()
	return testresults
}

func worker(jobChan <-chan test.TestTask, testresults *[]test.TestResult) {
	defer wg.Done()
	for job := range jobChan {
		fmt.Print("Enqueuing job: ")
		fmt.Println(job)
		wg.Add(1)
		go func(job test.TestTask, testresults *[]test.TestResult) {
			defer wg.Done()
			testResult, err := job.RunTest()
			if err != nil {
				fmt.Println(err)
			}
			*testresults = append(*testresults, testResult)
		}(job, testresults)
	}
}
