package test

import (
	"fmt"
	"strings"

	"github.com/johnhany97/grader/processors"
)

type OutputTestHandler struct {
	Test   Test
	File   string
	Folder string
}

func (opt OutputTestHandler) RunTest() (TestResult, error) {
	processor := processors.SubmissionsProcessor{}
	stdout, stderr, err := processor.Execute(opt.File, opt.Folder)
	if err != nil {
		return opt.handleErr(err, stdout, stderr)
	}
	return opt.NewResult(stdout, stderr), nil
}

func (opt OutputTestHandler) NewResult(stdout string, stderr string) TestResult {
	tr := TestResult{
		Test: opt.Test,
	}

	// Add Trimmed outputs
	tr.StdOut = strings.TrimSpace(stdout)
	tr.StdErr = strings.TrimSpace(stderr)

	// Add Programmatically generated description of test
	if tr.StdErr == "" {
		exp := strings.Join(opt.Test.ExpectedOutput, "\n")
		tr.Successful = tr.StdOut == exp
		tr.Description = fmt.Sprintf("Expected:\n%v\nGot:\n%v\nTest passed: %v", exp, tr.StdOut, tr.Successful)
	}

	return tr
}

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
