package test

import (
	"strings"

	"github.com/johnhany97/grader/processors"
)

// JavaStyleTestHandler is a struct containing all the properties
// needed to be able to execute a test task that does java style checking
type JavaStyleTestHandler struct {
	Test   Test   // The Test itself as extracted from the Schema
	File   string // Name of the file containing the code being assessed
	Folder string // Folder within which this file exists
}

// RunTest is a method used to run a test task
func (jst JavaStyleTestHandler) RunTest() (Result, error) {
	processor := processors.SubmissionsProcessor{}
	stdout, stderr, err := processor.ExecuteJavaStyle(jst.File, jst.Folder)
	if err != nil {
		return jst.handleErr(err, stdout, stderr)
	}
	return jst.NewResult(stdout, stderr), nil
}

// NewResult returns back the result of processing the output of
// executing the test task
func (jst JavaStyleTestHandler) NewResult(stdout string, stderr string) Result {
	tr := Result{
		Test: jst.Test,
	}

	// Add Trimmed outputs
	tr.StdOut = strings.TrimSpace(stdout)
	tr.StdErr = strings.TrimSpace(stderr)
	tr.Description = "Results generated by Checkstyle. A static analyzer built for analyzing Java Programs."

	return tr
}

// handleErr handles potnetial errors produced from the execution and
// customizes the test task output to reflect is the error was due to
// a certain specific cause (for example, a timeout error)
func (jst JavaStyleTestHandler) handleErr(e error, stdout string, stderr string) (Result, error) {
	if strings.Compare(e.Error(), "timeout") == 0 {
		return Result{
			Test:     jst.Test,
			TimedOut: true,
			StdOut:   strings.TrimSpace(stdout),
			StdErr:   strings.TrimSpace(stderr),
		}, nil
	}
	return Result{
		Test:   jst.Test,
		StdOut: strings.TrimSpace(stdout),
		StdErr: strings.TrimSpace(stderr),
	}, e
}
