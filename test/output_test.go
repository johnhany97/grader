package test

import (
	"errors"
	"fmt"
	"path/filepath"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestRunTestOutputTestHandler(t *testing.T) {
	className := "Solution"
	// Count the Ks
	folder := "testData/helloWorld/"
	availableFormats := []string{
		"java",
		"py",
		"cpp",
		"cs",
	}
	tests := []Test{
		Test{
			Type: "output",
			ExpectedOutput: []string{
				"Hello world",
			},
		},
	}

	var opts []OutputTestHandler
	var expectedResults []Result

	for _, test := range tests {
		for _, ext := range availableFormats {
			opts = append(opts, OutputTestHandler{
				Test:   test,
				File:   className + "." + ext,
				Folder: folder,
			})
			expectedResults = append(expectedResults, Result{
				Test:        test,
				StdOut:      strings.Join(test.ExpectedOutput, "\n"),
				StdErr:      "",
				Successful:  true,
				Similarity:  1,
				TimedOut:    false,
				Description: fmt.Sprintf("Expected:\n%v\nGot:\n%v\nTest passed: %v", strings.Join(test.ExpectedOutput, "\n"), strings.Join(test.ExpectedOutput, "\n"), true),
			})
		}
	}

	var actualResults []Result

	for _, opt := range opts {
		result, _ := opt.RunTest()
		actualResults = append(actualResults, result)
	}

	for i := 0; i < len(opts); i++ {
		correct := cmp.Equal(expectedResults[i], actualResults[i])
		if !correct {
			t.Fatalf("Failed\n- Language: %s\n- Expected result: \n%v\n- Actual result: \n%v\n", filepath.Ext(opts[i].File), expectedResults[i], actualResults[i])
		} else {
			t.Logf("Successfully graded the solution for %s language", filepath.Ext(opts[i].File))
		}
	}
}

func BenchmarkRunTestOutputTestHandler(b *testing.B) {
	opt := OutputTestHandler{
		Test: Test{
			Type: "output",
			ExpectedOutput: []string{
				"11",
				"3",
			},
		},
		File:   "Solution.java",
		Folder: "testData/helloWorld/",
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		opt.RunTest()
	}
}

func TestNewResultOutputTestHandler(t *testing.T) {
	test := Test{
		Type: "output",
		ExpectedOutput: []string{
			"11",
			"3",
		},
	}
	data :=
		[]struct {
			opt    OutputTestHandler
			Stdout string
			Stderr string
			TR     Result
		}{
			{
				opt: OutputTestHandler{
					Test: test,
				},
				Stdout: "11\n3",
				Stderr: "",
				TR: Result{
					Test:        test,
					StdOut:      "11\n3",
					StdErr:      "",
					Successful:  true,
					Similarity:  1,
					TimedOut:    false,
					Description: fmt.Sprintf("Expected:\n%v\nGot:\n%v\nTest passed: %v", strings.Join(test.ExpectedOutput, "\n"), strings.Join(test.ExpectedOutput, "\n"), true),
				},
			},
			{
				opt: OutputTestHandler{
					Test: test,
				},
				Stdout: "10\n2",
				Stderr: "",
				TR: Result{
					Test:        test,
					StdOut:      "10\n2",
					StdErr:      "",
					Successful:  false,
					Similarity:  0.6666666666666666,
					TimedOut:    false,
					Description: fmt.Sprintf("Expected:\n%v\nGot:\n%v\nTest passed: %v", strings.Join(test.ExpectedOutput, "\n"), "10\n2", false),
				},
			},
		}
	for _, entry := range data {
		tr := entry.opt.NewResult(entry.Stdout, entry.Stderr)
		correct := cmp.Equal(entry.TR, tr)
		if !correct {
			t.Fatalf("Failed\n- Expected result: \n%v\n- Actual result: \n%v\n", entry.TR, tr)
		} else {
			t.Logf("Successfully obtained the expected result")
		}
	}
}

func BenchmarkNewResultOutputTestHandler(b *testing.B) {
	test := Test{
		Type: "output",
		ExpectedOutput: []string{
			"11",
			"3",
		},
	}
	stdout := strings.Join(test.ExpectedOutput, "\n")
	stderr := ""
	opt := OutputTestHandler{
		Test:   test,
		File:   "Solution.java",
		Folder: "testData/helloWorld/",
	}
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		opt.NewResult(stdout, stderr)
	}
}

func TestHandleErrOutputTestHandler(t *testing.T) {
	test := Test{
		Type: "output",
		ExpectedOutput: []string{
			"Hello world",
		},
	}
	stdout := strings.Join(test.ExpectedOutput, "\n")
	stderr := ""
	opt := OutputTestHandler{
		Test:   test,
		File:   "Solution.java",
		Folder: "testData/helloWorld/",
	}
	err := errors.New("timeout")
	err2 := errors.New("something else")

	result1, _ := opt.handleErr(err, stdout, stderr)
	result2, _ := opt.handleErr(err2, stdout, stderr)

	correct1 := result1.StdOut == stdout && result1.StdErr == stderr && result1.TimedOut && !result1.Successful
	if !correct1 {
		t.Fatalf("Expected\n- StdOut: %s\n- StdErr: %s\n- Timeout: %v\n-------\nGot\n- StdOut: %s\n- StdErr: %s\n- Timeout: %v\n", stdout, stderr, true, result1.StdOut, result1.StdErr, result1.TimedOut)
	}
	correct2 := result2.StdOut == stdout && result2.StdErr == stderr && !result2.TimedOut && !result2.Successful
	if !correct2 {
		t.Fatalf("Expected\n- StdOut: %s\n- StdErr: %s\n- Timeout: %v\n-------\nGot\n- StdOut: %s\n- StdErr: %s\n- Timeout: %v\n", stdout, stderr, false, result2.StdOut, result2.StdErr, result2.TimedOut)
	}
}

func BenchmarkHandleErrOutputTestHandler(b *testing.B) {
	test := Test{
		Type: "output",
		ExpectedOutput: []string{
			"11",
			"3",
		},
	}
	stdout := strings.Join(test.ExpectedOutput, "\n")
	stderr := ""
	opt := OutputTestHandler{
		Test:   test,
		File:   "Solution.java",
		Folder: "testData/helloWorld/",
	}
	err := errors.New("timeout")

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		opt.handleErr(err, stdout, stderr)
	}
}
