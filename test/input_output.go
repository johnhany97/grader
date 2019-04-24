package test

import (
	"fmt"
	"strings"

	"github.com/johnhany97/grader/processors"
	"github.com/xrash/smetrics"
)

// InputOutputTestHandler is a struct containing all the properties
// needed to be able to execute a test task with input and expected output
// as the parameters.
type InputOutputTestHandler struct {
	Test   Test   // The Test itself as extracted from the Schema
	File   string // Name of the file containing the code being assessed
	Folder string // Folder within which this file exists
}

// RunTest is a method used to run a test task of type Input Output
func (iot InputOutputTestHandler) RunTest() (TestResult, error) {
	processor := processors.SubmissionsProcessor{}
	stdout, stderr, err := processor.ExecuteWithInput(iot.File, iot.Folder, strings.Join(iot.Test.Input, "\n"))
	if err != nil {
		return iot.handleErr(err)
	}
	return iot.NewResult(stdout, stderr), nil
}

// NewResult returns back the result of processing the output of
// executing the test task
func (iot InputOutputTestHandler) NewResult(stdout string, stderr string) TestResult {
	tr := TestResult{
		Test: iot.Test,
	}

	// Add Trimmed outputs
	tr.StdOut = strings.TrimSpace(stdout)
	tr.StdErr = strings.TrimSpace(stderr)

	successful := tr.StdOut == strings.Join(iot.Test.ExpectedOutput, "\n")
	tr.Successful = successful
	tr.Similarity = smetrics.JaroWinkler(tr.StdOut, strings.Join(iot.Test.ExpectedOutput, "\n"), 0.7, 4)

	// Add Programmatically generated description of test
	if tr.StdErr == "" {
		tr.Description = fmt.Sprintf("Expected:\n%v\nGot:\n%v\nTest passed: %v", strings.Join(iot.Test.ExpectedOutput, "\n"), tr.StdOut, successful)
	}

	return tr
}

// handleErr handles potnetial errors produced from the execution and
// customizes the test task output to reflect is the error was due to
// a certain specific cause (for example, a timeout error)
func (iot InputOutputTestHandler) handleErr(e error) (TestResult, error) {
	if strings.Compare(e.Error(), "timeout") == 0 {
		return TestResult{
			Test:     iot.Test,
			TimedOut: true,
		}, nil
	}
	return TestResult{
		Test: iot.Test,
	}, e
}
