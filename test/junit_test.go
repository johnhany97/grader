package test

import (
	"errors"
	"path/filepath"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestRunTestJUnitTestHandler(t *testing.T) {
	className := "Solution"
	// Count the Ks
	folder := "testData/adder/"
	availableFormats := []string{
		"java",
	}
	tests := []Test{
		Test{
			Type:     "junit",
			UnitTest: "@Test\n  public void adderWorksWithZero() {\n    Solution s = new Solution();\n    int actual = s.adder(0, 3);\n    assertEquals(3, actual);\n  }",
		},
		Test{
			Type:     "junit",
			UnitTest: "@Test\n  public void adderWorksWithNegative() {\n    Solution s = new Solution();\n    int actual = s.adder(0, -1);\n    assertEquals(-1, actual);\n  }",
		},
	}

	var juts []JUnitTestHandler
	var expectedResults []Result

	for _, test := range tests {
		for _, ext := range availableFormats {
			juts = append(juts, JUnitTestHandler{
				Test:      test,
				File:      className + "." + ext,
				Folder:    folder,
				ClassName: className,
			})
			expectedResults = append(expectedResults, Result{
				Test:        test,
				StdErr:      "",
				StdOut:      "JUnit version 4.10\r\n.\r\nTime: 0.009\r\n\r\nOK (1 test)",
				Successful:  true,
				Similarity:  0,
				TimedOut:    false,
				Description: "",
			})
		}
	}

	var actualResults []Result

	for _, jut := range juts {
		result, _ := jut.RunTest()
		actualResults = append(actualResults, result)
	}

	for i := 0; i < len(juts); i++ {
		correct := expectedResults[i].Successful == actualResults[i].Successful && actualResults[i].Successful && strings.Contains(actualResults[i].StdOut, "OK (1 test)")
		if !correct {
			t.Fatalf("Failed\n- Language: %s\n- Expected result: \n%v\n- Actual result: \n%v\n", filepath.Ext(juts[i].File), expectedResults[i], actualResults[i])
		} else {
			t.Logf("Successfully graded the solution for %s language", filepath.Ext(juts[i].File))
		}
	}
}

func BenchmarkRunTestJUnitTestHandler(b *testing.B) {
	jut := JUnitTestHandler{
		Test: Test{
			Type:     "junit",
			UnitTest: "@Test\n  public void adderWorksWithNegative() {\n    Solution s = new Solution();\n    int actual = s.adder(0, -1);\n    assertEquals(-1, actual);\n  }",
		},
		File:      "Solution.java",
		Folder:    "testData/adder/",
		ClassName: "Solution",
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		jut.RunTest()
	}
}

func TestNewResultJUnitTestHandler(t *testing.T) {
	test := Test{
		Type:     "junit",
		UnitTest: "@Test\n  public void adderWorksWithZero() {\n    Solution s = new Solution();\n    int actual = s.adder(0, 3);\n    assertEquals(3, actual);\n  }",
	}
	data := []struct {
		Jut    JUnitTestHandler
		Stdout string
		Stderr string
		TR     Result
	}{
		{
			Jut: JUnitTestHandler{
				Test: test,
			},
			Stdout: "JUnit version 4.10\r\n.\r\nTime: 0.009\r\n\r\nOK (1 test)",
			Stderr: "",
			TR: Result{
				Test:        test,
				StdOut:      "JUnit version 4.10\r\n.\r\nTime: 0.009\r\n\r\nOK (1 test)",
				StdErr:      "",
				Successful:  true,
				Similarity:  0,
				TimedOut:    false,
				Description: "",
			},
		},
		{
			Jut: JUnitTestHandler{
				Test: test,
			},
			Stdout: "JUnit version 4.10\r\n.\r\nTime: 0.009\r\n\r\nOK (1 test)",
			Stderr: "",
			TR: Result{
				Test:        test,
				StdOut:      "JUnit version 4.10\r\n.\r\nTime: 0.009\r\n\r\nOK (1 test)",
				StdErr:      "",
				Successful:  true,
				Similarity:  0,
				TimedOut:    false,
				Description: "",
			},
		},
	}
	for _, entry := range data {
		tr := entry.Jut.NewResult(entry.Stdout, entry.Stderr)
		correct := cmp.Equal(entry.TR, tr)
		if !correct {
			t.Fatalf("Failed\n- Actual result: \n%v\n", tr)
		} else {
			t.Logf("Successfully obtained the expected result")
		}
	}
}

func BenchmarkNewResultJUnitTestHandler(b *testing.B) {
	test := Test{
		Type:     "junit",
		UnitTest: "@Test\n  public void adderWorksWithNegative() {\n    Solution s = new Solution();\n    int actual = s.adder(0, -1);\n    assertEquals(-1, actual);\n  }",
	}
	stdout := "JUnit version 4.10\r\n.\r\nTime: 0.009\r\n\r\nOK (1 test)"
	stderr := ""
	jut := JUnitTestHandler{
		Test:      test,
		File:      "Solution.java",
		Folder:    "testData/adder/",
		ClassName: "Solution",
	}
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		jut.NewResult(stdout, stderr)
	}
}

func TestHandleErrJUnitTestHandler(t *testing.T) {
	test := Test{
		Type:     "junit",
		UnitTest: "@Test\n  public void adderWorksWithNegative() {\n    Solution s = new Solution();\n    int actual = s.adder(0, -1);\n    assertEquals(-1, actual);\n  }",
	}
	stdout := "JUnit version 4.10\r\n.\r\nTime: 0.009\r\n\r\nOK (1 test)"
	stderr := ""
	jut := JUnitTestHandler{
		Test:      test,
		File:      "Solution.java",
		Folder:    "testData/printer/",
		ClassName: "Solution",
	}
	err := errors.New("timeout")
	err2 := errors.New("something else")

	result1, _ := jut.handleErr(err, stdout)
	result2, _ := jut.handleErr(err2, stdout)

	correct1 := result1.StdOut == stdout && result1.StdErr == stderr && result1.TimedOut && !result1.Successful
	if !correct1 {
		t.Fatalf("Expected\n- StdOut: %s\n- StdErr: %s\n- Timeout: %v\n-------\nGot\n- StdOut: %s\n- StdErr: %s\n- Timeout: %v\n", stdout, stderr, true, result1.StdOut, result1.StdErr, result1.TimedOut)
	}
	correct2 := result2.StdOut == stdout && result2.StdErr == stderr && !result2.TimedOut && !result2.Successful
	if !correct2 {
		t.Fatalf("Expected\n- StdOut: %s\n- StdErr: %s\n- Timeout: %v\n-------\nGot\n- StdOut: %s\n- StdErr: %s\n- Timeout: %v\n", stdout, stderr, false, result2.StdOut, result2.StdErr, result2.TimedOut)
	}
}

func BenchmarkHandleErrJUnitTestHandler(b *testing.B) {
	test := Test{
		Type:     "junit",
		UnitTest: "@Test\n  public void adderWorksWithNegative() {\n    Solution s = new Solution();\n    int actual = s.adder(0, -1);\n    assertEquals(-1, actual);\n  }",
	}
	stdout := "JUnit version 4.10\r\n.\r\nTime: 0.009\r\n\r\nOK (1 test)"
	iot := JUnitTestHandler{
		Test:      test,
		File:      "Solution.java",
		Folder:    "testData/printer/",
		ClassName: "Solution",
	}
	err := errors.New("timeout")

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		iot.handleErr(err, stdout)
	}
}
