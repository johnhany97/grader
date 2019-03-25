package test

import (
	"fmt"
	"strings"

	"github.com/johnhany97/grader/processors"
)

type InputOutputTestHandler struct {
	Test   Test
	File   string
	Folder string
}

func (iot InputOutputTestHandler) RunTest() (TestResult, error) {
	processor := processors.Processor{}
	stdout, stderr := processor.ExecuteWithInput(iot.File, iot.Folder, strings.Join(iot.Test.Input, "\n"))
	return iot.NewResult(stdout, stderr), nil
}

func (iot InputOutputTestHandler) NewResult(stdout string, stderr string) TestResult {
	tr := TestResult{
		Test: iot.Test,
	}

	// Add Trimmed outputs
	tr.StdOut = strings.TrimSpace(stdout)
	tr.StdErr = strings.TrimSpace(stderr)

	// Add Programmatically generated description of test
	if tr.StdErr == "" {
		tr.Description = fmt.Sprintf("Expected:\n%v\nGot:\n%v\nTest passed: %v", strings.Join(iot.Test.ExpectedOutput, "\n"), tr.StdOut, tr.StdOut == strings.Join(iot.Test.ExpectedOutput, "\n"))
	}

	return tr
}
