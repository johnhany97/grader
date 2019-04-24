package test

import (
	"fmt"
	"strings"

	"github.com/johnhany97/grader/processors"
	"github.com/xrash/smetrics"
)

// OutputTestHandler is a struct containing all the properties
// needed to be able to execute a test task with just the expected output
// as the parameter.
type OutputTestHandler struct {
	Test   Test   // The Test itself as extracted from the Schema
	File   string // Name of the file containing the code being assessed
	Folder string // Folder within which this file exists
}

// RunTest is a method used to run a test task
func (opt OutputTestHandler) RunTest() (TestResult, error) {
	processor := processors.SubmissionsProcessor{}
	stdout, stderr, err := processor.Execute(opt.File, opt.Folder)
	if err != nil {
		return opt.handleErr(err, stdout, stderr)
	}
	return opt.NewResult(stdout, stderr), nil
}

// NewResult returns back the result of processing the output of
// executing the test task
func (opt OutputTestHandler) NewResult(stdout string, stderr string) TestResult {
	tr := TestResult{
		Test: opt.Test,
	}

	// Add Trimmed outputs
	tr.StdOut = strings.TrimSpace(stdout)
	tr.StdErr = strings.TrimSpace(stderr)
	tr.Similarity = smetrics.JaroWinkler(tr.StdOut, strings.Join(opt.Test.ExpectedOutput, "\n"), 0.7, 4)

	// Add Programmatically generated description of test
	if tr.StdErr == "" {
		exp := strings.Join(opt.Test.ExpectedOutput, "\n")
		tr.Successful = tr.StdOut == exp
		tr.Description = fmt.Sprintf("Expected:\n%v\nGot:\n%v\nTest passed: %v", exp, tr.StdOut, tr.Successful)
	}

	return tr
}

// handleErr handles potnetial errors produced from the execution and
// customizes the test task output to reflect is the error was due to
// a certain specific cause (for example, a timeout error)
func (opt OutputTestHandler) handleErr(e error, stdout string, stderr string) (TestResult, error) {
	if strings.Compare(e.Error(), "timeout") == 0 {
		return TestResult{
			Test:     opt.Test,
			TimedOut: true,
			StdOut:   strings.TrimSpace(stdout),
			StdErr:   strings.TrimSpace(stderr),
		}, nil
	}
	return TestResult{
		Test:   opt.Test,
		StdOut: strings.TrimSpace(stdout),
		StdErr: strings.TrimSpace(stderr),
	}, e
}
