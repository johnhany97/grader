package test

import (
	"errors"
	"path/filepath"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestRunTestPyUnitTestHandler(t *testing.T) {
	className := "Solution"
	folder := "testData/odd/"
	availableFormats := []string{
		"py",
	}
	tests := []Test{
		Test{
			Type:     "pyunit",
			UnitTest: "\tdef test_is_odd_on_odd(self):\n\t\tself.assertTrue(is_odd(5))",
		},
		Test{
			Type:     "pyunit",
			UnitTest: "\tdef test_is_odd_on_even(self):\n\t\tself.assertFalse(is_odd(2))",
		},
	}

	var puts []PyUnitTestHandler

	for _, test := range tests {
		for _, ext := range availableFormats {
			puts = append(puts, PyUnitTestHandler{
				Test:      test,
				File:      className + "." + ext,
				Folder:    folder,
				ClassName: className,
			})
		}
	}

	var actualResults []Result

	for _, put := range puts {
		result, _ := put.RunTest()
		actualResults = append(actualResults, result)
	}

	for i := 0; i < len(puts); i++ {
		correct := actualResults[i].Successful && strings.Contains(actualResults[i].StdOut, "OK")
		if !correct {
			t.Fatalf("Failed\n- Actual result: \n%v\n", actualResults[i])
		} else {
			t.Logf("Successfully graded the solution for %s language", filepath.Ext(puts[i].File))
		}
	}
}

func TestNewResultPyUnitTestHandler(t *testing.T) {
	test := Test{
		Type:     "pyunit",
		UnitTest: "\tdef test_is_odd_on_odd(self):\n\t\tself.assertTrue(is_odd(5))",
	}
	data := []struct {
		Put    PyUnitTestHandler
		Stdout string
		Stderr string
		TR     Result
	}{
		{
			Put: PyUnitTestHandler{
				Test: test,
			},
			Stdout: ".\n----------------------------------------------------------------------\nRan 1 test in 0.000s\n\nOK",
			Stderr: "",
			TR: Result{
				Test:        test,
				StdOut:      ".\n----------------------------------------------------------------------\nRan 1 test in 0.000s\n\nOK",
				StdErr:      "",
				Successful:  true,
				Similarity:  0,
				TimedOut:    false,
				Description: "",
			},
		},
		{
			Put: PyUnitTestHandler{
				Test: test,
			},
			Stdout: "",
			Stderr: "AssertionError: False is not true",
			TR: Result{
				Test:        test,
				StdOut:      "",
				StdErr:      "AssertionError: False is not true",
				Successful:  false,
				Similarity:  0,
				TimedOut:    false,
				Description: "",
			},
		},
	}
	for _, entry := range data {
		tr := entry.Put.NewResult(entry.Stdout, entry.Stderr)
		correct := cmp.Equal(entry.TR, tr)
		if !correct {
			t.Fatalf("Failed\n- Actual result: \n%v\n", tr)
		} else {
			t.Logf("Successfully obtained the expected result")
		}
	}
}

func BenchmarkNewResultPyUnitTestHandler(b *testing.B) {
	test := Test{
		Type:     "pyunit",
		UnitTest: "\tdef test_is_odd_on_odd(self):\n\t\tself.assertTrue(is_odd(5))",
	}
	stdout := ".\n----------------------------------------------------------------------\nRan 1 test in 0.000s\n\nOK"
	stderr := ""
	put := PyUnitTestHandler{
		Test:      test,
		File:      "Solution.py",
		Folder:    "testData/odd/",
		ClassName: "Solution",
	}
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		put.NewResult(stdout, stderr)
	}
}

func TestHandleErrPyUnitTestHandler(t *testing.T) {
	test := Test{
		Type:     "pyunit",
		UnitTest: "\tdef test_is_odd_on_odd(self):\n\t\tself.assertTrue(is_odd(5))",
	}
	stdout := ".\n----------------------------------------------------------------------\nRan 1 test in 0.000s\n\nOK"
	stderr := ""
	put := PyUnitTestHandler{
		Test:      test,
		File:      "Solution.py",
		Folder:    "testData/odd/",
		ClassName: "Solution",
	}
	err := errors.New("timeout")
	err2 := errors.New("something else")

	result1, _ := put.handleErr(err, stdout)
	result2, _ := put.handleErr(err2, stdout)

	correct1 := result1.StdOut == stdout && result1.StdErr == stderr && result1.TimedOut && !result1.Successful
	if !correct1 {
		t.Fatalf("Expected\n- StdOut: %s\n- StdErr: %s\n- Timeout: %v\n-------\nGot\n- StdOut: %s\n- StdErr: %s\n- Timeout: %v\n", stdout, stderr, true, result1.StdOut, result1.StdErr, result1.TimedOut)
	}
	correct2 := result2.StdOut == stdout && result2.StdErr == stderr && !result2.TimedOut && !result2.Successful
	if !correct2 {
		t.Fatalf("Expected\n- StdOut: %s\n- StdErr: %s\n- Timeout: %v\n-------\nGot\n- StdOut: %s\n- StdErr: %s\n- Timeout: %v\n", stdout, stderr, false, result2.StdOut, result2.StdErr, result2.TimedOut)
	}
}

func BenchmarkHandleErrPyUnitTestHandler(b *testing.B) {
	test := Test{
		Type:     "pyunit",
		UnitTest: "\tdef test_is_odd_on_odd(self):\n\t\tself.assertTrue(is_odd(5))",
	}
	stdout := ".\n----------------------------------------------------------------------\nRan 1 test in 0.000s\n\nOK"
	iot := PyUnitTestHandler{
		Test:      test,
		File:      "Solution.py",
		Folder:    "testData/odd/",
		ClassName: "Solution",
	}
	err := errors.New("timeout")

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		iot.handleErr(err, stdout)
	}
}
