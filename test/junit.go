package test

import (
	"fmt"
	"io/ioutil"
	"strings"

	"github.com/johnhany97/grader/processors"
)

type JUnitTestHandler struct {
	Test      Test
	File      string
	Folder    string
	ClassName string
}

func (jut JUnitTestHandler) RunTest() (TestResult, error) {
	processor := processors.Processor{}
	// Obtain Junit file shell
	junitShell, err := ioutil.ReadFile("assets/JUnit.java")
	if err != nil {
		handleErr(err)
	}
	// Obtain all unit tests
	// Put it all together
	junitFinal := fmt.Sprintf(string(junitShell), jut.ClassName, jut.Test.UnitTest)
	stdout, err := processor.ExecuteJUnitTests(jut.ClassName, jut.Folder, junitFinal)
	if err != nil {
		handleErr(err)
	}
	return jut.NewResult(stdout, ""), nil
}

func (jut JUnitTestHandler) NewResult(stdout string, stderr string) TestResult {
	tr := TestResult{
		Test: jut.Test,
	}

	// Add Trimmed outputs
	tr.StdOut = strings.TrimSpace(stdout)
	tr.StdErr = strings.TrimSpace(stderr)

	return tr
}

func handleErr(e error) (TestResult, error) {
	return TestResult{}, e
}
