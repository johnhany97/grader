package test

import (
	"strings"

	"github.com/johnhany97/grader/processors"
)

type JavaStyleTestHandler struct {
	Test   Test
	File   string
	Folder string
}

func (jst JavaStyleTestHandler) RunTest() (TestResult, error) {
	processor := processors.SubmissionsProcessor{}
	stdout, stderr := processor.ExecuteJavaStyle(jst.File, jst.Folder)
	return jst.NewResult(stdout, stderr), nil
}

func (jst JavaStyleTestHandler) NewResult(stdout string, stderr string) TestResult {
	tr := TestResult{
		Test: jst.Test,
	}

	// Add Trimmed outputs
	tr.StdOut = strings.TrimSpace(stdout)
	tr.StdErr = strings.TrimSpace(stderr)
	tr.Description = "Results generated by Checkstyle. A static analyzer built for analyzing Java Programs."

	return tr
}
