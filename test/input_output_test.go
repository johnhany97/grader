package test

import (
	"fmt"
	"path/filepath"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestRunTestInputOutputTestHandler(t *testing.T) {
	className := "Solution"
	// Count the Ks
	folder := "testData/printer/"
	availableFormats := []string{
		"java",
		"py",
		"cpp",
		"cs",
	}
	tests := []Test{
		Test{
			Type: "io",
			Input: []string{
				"3",
				"2",
				"4",
				"5",
			},
			ExpectedOutput: []string{
				"3",
				"5",
				"6",
			},
		},
		Test{
			Type: "io",
			Input: []string{
				"2",
				"1",
				"2",
			},
			ExpectedOutput: []string{
				"2",
				"3",
			},
		},
		Test{
			Type: "io",
			Input: []string{
				"2",
				"10",
				"2",
			},
			ExpectedOutput: []string{
				"11",
				"3",
			},
		},
	}

	var iots []InputOutputTestHandler
	var expectedResults []Result

	for _, test := range tests {
		for _, ext := range availableFormats {
			iots = append(iots, InputOutputTestHandler{
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

	for _, iot := range iots {
		result, _ := iot.RunTest()
		actualResults = append(actualResults, result)
	}

	for i := 0; i < len(iots); i++ {
		correct := cmp.Equal(expectedResults[i], actualResults[i])
		if !correct {
			t.Fatalf("Failed\n- Language: %s\n- Expected result: \n%v\n- Actual result: \n%v\n", filepath.Ext(iots[i].File), expectedResults[i], actualResults[i])
		} else {
			t.Logf("Successfully graded the solution for %s language", filepath.Ext(iots[i].File))
		}
	}
}

func TestNewResultInputOutputTestHandler(t *testing.T) {
	test := Test{
		Type: "io",
		Input: []string{
			"2",
			"10",
			"2",
		},
		ExpectedOutput: []string{
			"11",
			"3",
		},
	}
	data :=
		[]struct {
			Iot    InputOutputTestHandler
			Stdout string
			Stderr string
			TR     Result
		}{
			{
				Iot: InputOutputTestHandler{
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
				Iot: InputOutputTestHandler{
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
		tr := entry.Iot.NewResult(entry.Stdout, entry.Stderr)
		correct := cmp.Equal(entry.TR, tr)
		if !correct {
			t.Fatalf("Failed\n- Expected result: \n%v\n- Actual result: \n%v\n", entry.TR, tr)
		} else {
			t.Logf("Successfully obtained the expected result")
		}
	}
}
