package test

import (
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

	var actualResults []TestResult

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
		TR     TestResult
	}{
		{
			Put: PyUnitTestHandler{
				Test: test,
			},
			Stdout: "Ran Tests.... OK",
			Stderr: "",
			TR: TestResult{
				Test:        test,
				StdOut:      "Ran Tests.... OK",
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
			TR: TestResult{
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