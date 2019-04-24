package grader

import (
	"sync"

	"github.com/johnhany97/grader/test"
)

// Schema defining all the parameters presented within
// a grading marking schema
type Schema struct {
	File      string      `json:"file"`      // Name of the file that is to be grader - ex. Solution.java
	Language  string      `json:"language"`  // One of the supported languages - check the inner supported languages
	ClassName string      `json:"className"` // ClassName is the name of the class in the product.
	Folder    string      `json:"folder"`    // Folder within which the file is contained
	Tests     []test.Test `json:"tests"`     // List of tests which actually define the marking schema
	Outfile   string      `json:"outfile"`   // Empty to output to terminal or a file name to which we store the results within the previously provided folder
}

// Grader defining required parameters to be presented
// to the grader to start the grading process
type Grader struct {
	Schema   Schema `json:"schema"`   // The schema as previously defined
	MaxTasks int    `jsno:"maxTasks"` // Limit to concurrent number of tasks executed at the same time
}

var wg sync.WaitGroup

// NewGrader is used to create a new Grader instance
func NewGrader(s Schema, mt int) Grader {
	return Grader{
		Schema:   s,
		MaxTasks: mt,
	}
}

// Grade is used to start the grading process over the Grader's
// properties by starting a job worker executing the tasks till
// they're all done
func (g Grader) Grade() []test.Result {
	// make a channel with a capacity of 100.
	jobChan := make(chan test.Task, g.MaxTasks)

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

// worker which given the channel and the test results
// slice awaits the tasks to be finished and appends them
// to the slice then signals to the waitGroup that we're done
func worker(jobChan <-chan test.Task, testresults *[]test.Result) {
	defer wg.Done()
	for job := range jobChan {
		wg.Add(1)
		go func(job test.Task, testresults *[]test.Result) {
			defer wg.Done()
			testResult, _ := job.RunTest()
			*testresults = append(*testresults, testResult)
		}(job, testresults)
	}
}
