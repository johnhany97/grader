package test

import (
	"fmt"
	"strings"

	"github.com/johnhany97/grader/processors"
)

// JUnitTestHandler is a struct containing all the properties
// needed to be able to execute a test task given a java unit
// test (junit) function and the file being tested as the parameters.
type JUnitTestHandler struct {
	Test      Test   // The Test itself as extracted from the Schema
	File      string // Name of the file containing the code being assessed
	Folder    string // Folder within which this file exists
	ClassName string // Name of the class
}

// RunTest is a method used to run a test task
func (jut JUnitTestHandler) RunTest() (Result, error) {
	processor := processors.SubmissionsProcessor{}
	// Obtain all unit tests
	// Put it all together
	junitFinal := fmt.Sprintf(string(junitShell), jut.ClassName, jut.Test.UnitTest)
	stdout, err := processor.ExecuteJUnitTests(jut.ClassName, jut.Folder, junitFinal)
	if err != nil {
		return jut.handleErr(err, stdout)
	}
	return jut.NewResult(stdout, ""), nil
}

// NewResult returns back the result of processing the output of
// executing the test task
func (jut JUnitTestHandler) NewResult(stdout string, stderr string) Result {
	tr := Result{
		Test: jut.Test,
	}

	// Add Trimmed outputs
	tr.StdOut = strings.TrimSpace(stdout)
	tr.StdErr = strings.TrimSpace(stderr)

	tr.Successful = strings.Contains(stdout, "OK (1 test)")

	return tr
}

// handleErr handles potnetial errors produced from the execution and
// customizes the test task output to reflect is the error was due to
// a certain specific cause (for example, a timeout error)
func (jut JUnitTestHandler) handleErr(e error, stdout string) (Result, error) {
	if strings.Compare(e.Error(), "timeout") == 0 {
		return Result{
			Test:     jut.Test,
			TimedOut: true,
			StdOut:   strings.TrimSpace(stdout),
		}, nil
	}
	return Result{
		Test:   jut.Test,
		StdOut: strings.TrimSpace(stdout),
	}, e
}

const junitShell = `import java.util.*;

import org.junit.Test;
import org.junit.runner.JUnitCore;
import org.junit.runner.Result;
import org.junit.runner.notification.Failure;

import static org.junit.Assert.*;

public class %vTest {
	%v
}`
