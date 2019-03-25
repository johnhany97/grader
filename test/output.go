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
	processor := processors.Processor{}
	stdout, stderr := processor.Execute(opt.File, opt.Folder)
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
		tr.Description = fmt.Sprintf("Expected:\n%v\nGot:\n%v\nTest passed: %v", exp, tr.StdOut, tr.StdOut == exp)
	}

	return tr
}
