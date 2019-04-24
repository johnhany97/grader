package test

import (
	"errors"
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

func BenchmarkRunTestJavaStyleTestHandler(b *testing.B) {
	jst := JavaStyleTestHandler{
		Test: Test{
			Type: "javaStyle",
		},
		File:   "Solution.java",
		Folder: "testData/printer/",
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		jst.RunTest()
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

func BenchmarkNewResultJavaStyleTestHandler(b *testing.B) {
	test := Test{
		Type: "javaStyle",
	}
	stdout := "[Audit Started]\n[Audit Done]"
	stderr := ""
	jst := JavaStyleTestHandler{
		Test:   test,
		File:   "Solution.java",
		Folder: "testData/printer/",
	}
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		jst.NewResult(stdout, stderr)
	}
}

func TestHandleErrJavaStyleTestHandler(t *testing.T) {
	test := Test{
		Type: "javaStyle",
	}
	stdout := "[Audit Started]\n[Audit Done]"
	stderr := ""
	iot := JavaStyleTestHandler{
		Test:   test,
		File:   "Solution.java",
		Folder: "testData/printer/",
	}
	err := errors.New("timeout")
	err2 := errors.New("something else")

	result1, _ := iot.handleErr(err, stdout, stderr)
	result2, _ := iot.handleErr(err2, stdout, stderr)

	correct1 := result1.StdOut == stdout && result1.StdErr == stderr && result1.TimedOut && !result1.Successful
	if !correct1 {
		t.Fatalf("Expected\n- StdOut: %s\n- StdErr: %s\n- Timeout: %v\n-------\nGot\n- StdOut: %s\n- StdErr: %s\n- Timeout: %v\n", stdout, stderr, true, result1.StdOut, result1.StdErr, result1.TimedOut)
	}
	correct2 := result2.StdOut == stdout && result2.StdErr == stderr && !result2.TimedOut && !result2.Successful
	if !correct2 {
		t.Fatalf("Expected\n- StdOut: %s\n- StdErr: %s\n- Timeout: %v\n-------\nGot\n- StdOut: %s\n- StdErr: %s\n- Timeout: %v\n", stdout, stderr, false, result2.StdOut, result2.StdErr, result2.TimedOut)
	}
}

func BenchmarkHandleErrJavaStyleTestHandler(b *testing.B) {
	test := Test{
		Type: "javaStyle",
	}
	stdout := "[Audit Started]\n[Audit Done]"
	stderr := ""
	iot := JavaStyleTestHandler{
		Test:   test,
		File:   "Solution.java",
		Folder: "testData/printer/",
	}
	err := errors.New("timeout")

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		iot.handleErr(err, stdout, stderr)
	}
}
