package test

import (
	"fmt"
	"strings"

	"github.com/johnhany97/grader/processors"
)

type PyUnitTestHandler struct {
	Test      Test
	File      string
	Folder    string
	ClassName string
}

func (put PyUnitTestHandler) RunTest() (TestResult, error) {
	processor := processors.SubmissionsProcessor{}
	// Obtain all unit tests
	pyUnitFinal := fmt.Sprintf(string(pyUnitShell), put.ClassName, put.ClassName, put.Test.UnitTest)
	stdout, err := processor.ExecutePyUnitTests(put.File, put.ClassName, put.Folder, pyUnitFinal)
	if err != nil {
		return put.handleErr(err)
	}
	return put.NewResult(stdout, ""), nil
}

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

func (put PyUnitTestHandler) handleErr(e error) (TestResult, error) {
	if strings.Compare(e.Error(), "timeout") == 0 {
		return TestResult{
			Test:     put.Test,
			TimedOut: true,
		}, nil
	}
	return TestResult{
		Test: put.Test,
	}, e
}

const pyUnitShell = `
import unittest
from %v import *

class %vTestCase(unittest.TestCase):
%v

if __name__ == '__main__':
	unittest.main()`
