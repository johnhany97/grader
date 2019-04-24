package test

import (
	"path/filepath"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestRunTestJavaStyleTestHandler(t *testing.T) {
	className := "Solution"
	// Count the Ks
	folder := "testData/style/"
	availableFormats := []string{
		"java",
	}
	tests := []Test{
		Test{
			Type: "javaStyle",
		},
	}

	var jsts []JavaStyleTestHandler

	for _, test := range tests {
		for _, ext := range availableFormats {
			jsts = append(jsts, JavaStyleTestHandler{
				Test:   test,
				File:   className + "." + ext,
				Folder: folder,
			})
		}
	}

	var actualResults []Result

	for _, jst := range jsts {
		result, _ := jst.RunTest()
		actualResults = append(actualResults, result)
	}

	for i := 0; i < len(jsts); i++ {
		correct := strings.Contains(actualResults[i].StdOut, "Using the '.*' form of import should be avoided - java.util.*. [AvoidStarImport]")
		if !correct {
			t.Fatalf("Failed\n- Actual result: \n%v\n", actualResults[i])
		} else {
			t.Logf("Successfully style checked the solution for %s language", filepath.Ext(jsts[i].File))
		}
	}
}

func TestNewResultJavaStyleTestHandler(t *testing.T) {
	test := Test{
		Type: "javaStyle",
	}
	data := []struct {
		jst    JavaStyleTestHandler
		Stdout string
		Stderr string
		TR     Result
	}{
		{
			jst: JavaStyleTestHandler{
				Test: test,
			},
			Stdout: "The result of the audit",
			Stderr: "",
			TR: Result{
				Test:        test,
				StdOut:      "The result of the audit",
				StdErr:      "",
				Successful:  false,
				Similarity:  0,
				TimedOut:    false,
				Description: "Results generated by Checkstyle. A static analyzer built for analyzing Java Programs.",
			},
		},
		{
			jst: JavaStyleTestHandler{
				Test: test,
			},
			Stdout: "",
			Stderr: "Failed to start",
			TR: Result{
				Test:        test,
				StdOut:      "",
				StdErr:      "Failed to start",
				Successful:  false,
				Similarity:  0,
				TimedOut:    false,
				Description: "Results generated by Checkstyle. A static analyzer built for analyzing Java Programs.",
			},
		},
	}
	for _, entry := range data {
		tr := entry.jst.NewResult(entry.Stdout, entry.Stderr)
		correct := cmp.Equal(entry.TR, tr)
		if !correct {
			t.Fatalf("Failed\n- Actual result: \n%v\n", tr)
		} else {
			t.Logf("Successfully obtained the expected result")
		}
	}
}
