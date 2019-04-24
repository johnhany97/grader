package grader

import (
	"sync"

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
	Schema   Schema `json:"schema"`
	MaxTasks int    `jsno:"maxTasks"`
}

var wg sync.WaitGroup

func NewGrader(s Schema, mt int) Grader {
	return Grader{
		Schema:   s,
		MaxTasks: mt,
	}
}

func (g Grader) Grade() []test.Result {
	// make a channel with a capacity of 100.
	jobChan := make(chan test.Task, g.MaxTasks)

	// results
	testresults := []test.Result{}
	// start the worker
	wg.Add(1)
	go worker(jobChan, &testresults)

	for i := 0; i < len(g.Schema.Tests); i++ {
		// get test
		t := g.Schema.Tests[i]
		// create test task
		var task test.Task
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
		case "pyunit":
			task = test.PyUnitTestHandler{
				Test:      t,
				File:      g.Schema.File,
				Folder:    g.Schema.Folder,
				ClassName: g.Schema.ClassName,
			}
		case "javaStyle":
			task = test.JavaStyleTestHandler{
				Test:   t,
				File:   g.Schema.File,
				Folder: g.Schema.Folder,
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

func worker(jobChan <-chan test.Task, testresults *[]test.Result) {
	defer wg.Done()
	for job := range jobChan {
		// fmt.Print("Enqueuing job: ")
		// fmt.Println(job)
		wg.Add(1)
		go func(job test.Task, testresults *[]test.Result) {
			defer wg.Done()
			testResult, _ := job.RunTest()
			// if err != nil {
			// 	fmt.Println(err)
			// }
			*testresults = append(*testresults, testResult)
		}(job, testresults)
	}
}
