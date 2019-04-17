package test

import (
	"fmt"
	"strings"

	"github.com/johnhany97/grader/processors"
	"github.com/xrash/smetrics"
)

type InputOutputTestHandler struct {
	Test   Test
	File   string
	Folder string
}

func (iot InputOutputTestHandler) RunTest() (TestResult, error) {
	processor := processors.SubmissionsProcessor{}
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

	successful := tr.StdOut == strings.Join(iot.Test.ExpectedOutput, "\n")
	tr.Successful = successful
	tr.Similarity = smetrics.JaroWinkler(tr.StdOut, strings.Join(iot.Test.ExpectedOutput, "\n"), 0.7, 4)

	// Add Programmatically generated description of test
	if tr.StdErr == "" {
		tr.Description = fmt.Sprintf("Expected:\n%v\nGot:\n%v\nTest passed: %v", strings.Join(iot.Test.ExpectedOutput, "\n"), tr.StdOut, successful)
	}

	return tr
}
