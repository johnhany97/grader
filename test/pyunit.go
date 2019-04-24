package test

import (
	"fmt"
	"strings"

	"github.com/johnhany97/grader/processors"
)

// PyUnitTestHandler is a struct containing all the properties
// needed to be able to execute a test task given a python unit
// test function and te file being tested as the parameters.
type PyUnitTestHandler struct {
	Test      Test
	File      string
	Folder    string
	ClassName string
}

// RunTest is a method used to run a test task
func (put PyUnitTestHandler) RunTest() (TestResult, error) {
	processor := processors.SubmissionsProcessor{}
	// Obtain all unit tests
	pyUnitFinal := fmt.Sprintf(string(pyUnitShell), put.ClassName, put.ClassName, put.Test.UnitTest)
	stdout, err := processor.ExecutePyUnitTests(put.File, put.ClassName, put.Folder, pyUnitFinal)
	if err != nil {
		return put.handleErr(err, stdout)
	}
	return put.NewResult(stdout, ""), nil
}

// NewResult returns back the result of processing the output of
// executing the test task
func (put PyUnitTestHandler) NewResult(stdout string, stderr string) TestResult {
	tr := TestResult{
		Test: put.Test,
	}
	// Add Trimmed outputs
	tr.StdOut = strings.TrimSpace(stdout)
	tr.StdErr = strings.TrimSpace(stderr)

	tr.Successful = strings.Contains(stdout, "OK")

	return tr
}

func (put PyUnitTestHandler) handleErr(e error, stdout string) (TestResult, error) {
	if strings.Compare(e.Error(), "timeout") == 0 {
		return TestResult{
			Test:     put.Test,
			TimedOut: true,
			StdOut:   strings.TrimSpace(stdout),
		}, nil
	}
	return TestResult{
		Test:   put.Test,
		StdOut: strings.TrimSpace(stdout),
	}, e
}

const pyUnitShell = `
import unittest
from %v import *

class %vTestCase(unittest.TestCase):
%v

if __name__ == '__main__':
	unittest.main()`
