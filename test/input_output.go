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
	stdout, stderr, err := processor.ExecuteWithInput(iot.File, iot.Folder, strings.Join(iot.Test.Input, "\n"))
	if err != nil {
		return iot.handleErr(err, stdout, stderr)
	}
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

func (iot InputOutputTestHandler) handleErr(e error, stdout string, stderr string) (TestResult, error) {
	if strings.Compare(e.Error(), "timeout") == 0 {
		return TestResult{
			Test:     iot.Test,
			TimedOut: true,
			StdOut:   strings.TrimSpace(stdout),
			StdErr:   strings.TrimSpace(stderr),
		}, nil
	}
	return TestResult{
		Test:   iot.Test,
		StdOut: strings.TrimSpace(stdout),
		StdErr: strings.TrimSpace(stderr),
	}, e
}
