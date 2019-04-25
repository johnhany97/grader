package processors

import (
	"fmt"
	"strings"
	"testing"
)

func TestExecute(t *testing.T) {
	p := SubmissionsProcessor{}
	folder := "testData/helloWorld/"
	stdout := "Hello world"
	stderr := ""
	languages := []string{
		"cpp",
		"cs",
		"java",
		"py",
	}
	for _, lang := range languages {
		file := "Solution." + lang
		sout, serr, err := p.Execute(file, folder)
		if err != nil || strings.Compare(strings.TrimSpace(serr), stderr) != 0 || strings.Compare(strings.TrimSpace(sout), stdout) != 0 {
			t.Fatalf("Err should be %v, got %v\n StdOut should be %v got %v\n StdErr should be %v got %v", nil, err, stdout, sout, stderr, serr)
		}
	}
}

func BenchmarkExecute(b *testing.B) {
	p := SubmissionsProcessor{}
	file := "Solution.java"
	folder := "testData/helloWorld/"

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		p.Execute(file, folder)
	}
}

func TestExecuteWithInput(t *testing.T) {
	p := SubmissionsProcessor{}
	folder := "testData/printer/"
	stdin := "3\n4\n1\n3"
	stdout := "5\n2\n4"
	stderr := ""
	languages := []string{
		"cpp",
		"cs",
		"java",
		"py",
	}
	for _, lang := range languages {
		file := "Solution." + lang
		sout, serr, err := p.ExecuteWithInput(file, folder, stdin)
		if err != nil || strings.Compare(strings.TrimSpace(serr), stderr) != 0 || strings.Compare(strings.TrimSpace(sout), stdout) != 0 {
			t.Fatalf("Err should be %v, got %v\n StdOut should be %v got %v\n StdErr should be %v got %v", nil, err, stdout, sout, stderr, serr)
		}
	}
}

func BenchmarkExecuteWithInput(b *testing.B) {
	p := SubmissionsProcessor{}
	file := "Solution.java"
	folder := "testData/printer/"
	stdin := "3\n4\n1\n3"

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		p.ExecuteWithInput(file, folder, stdin)
	}
}

func TestExecuteJUnitTests(t *testing.T) {
	p := SubmissionsProcessor{}
	className := "Solution"
	folder := "testData/adder/"
	unitTest := "@Test\n  public void adderWorksWithNegative() {\n    Solution s = new Solution();\n    int actual = s.adder(0, -1);\n    assertEquals(-1, actual);\n  }"

	junitShell := `import java.util.*;

import org.junit.Test;
import org.junit.runner.JUnitCore;
import org.junit.runner.Result;
import org.junit.runner.notification.Failure;

import static org.junit.Assert.*;

public class %vTest {
	%v
}`

	junitFinal := fmt.Sprintf(string(junitShell), className, unitTest)

	sout, err := p.ExecuteJUnitTests(className, folder, junitFinal)
	if err != nil || !strings.Contains(sout, "OK") {
		t.Fatalf("Err should be %v, got %v\n StdOut should be contain \"OK\", got %v", nil, err, sout)
	}
}

func BenchmarkExecuteJUnitTests(b *testing.B) {
	p := SubmissionsProcessor{}
	className := "Solution"
	folder := "testData/adder/"
	unitTest := "@Test\n  public void adderWorksWithNegative() {\n    Solution s = new Solution();\n    int actual = s.adder(0, -1);\n    assertEquals(-1, actual);\n  }"

	junitShell := `import java.util.*;

import org.junit.Test;
import org.junit.runner.JUnitCore;
import org.junit.runner.Result;
import org.junit.runner.notification.Failure;

import static org.junit.Assert.*;

public class %vTest {
	%v
}`

	junitFinal := fmt.Sprintf(string(junitShell), className, unitTest)

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		p.ExecuteJUnitTests(className, folder, junitFinal)
	}
}

func TestExecutePyUnitTests(t *testing.T) {
	p := SubmissionsProcessor{}
	className := "Solution"
	file := className + ".py"
	folder := "testData/odd/"
	unitTest := "\tdef test_is_odd_on_odd(self):\n\t\tself.assertTrue(is_odd(5))"

	pyunitShell := `
import unittest
from %v import *

class %vTestCase(unittest.TestCase):
%v

if __name__ == '__main__':
	unittest.main()`

	pyunitFinal := fmt.Sprintf(string(pyunitShell), className, className, unitTest)

	sout, err := p.ExecutePyUnitTests(file, className, folder, pyunitFinal)
	if err != nil || !strings.Contains(sout, "OK") {
		t.Fatalf("Err should be %v, got %v\n StdOut should be contain \"OK\", got %v", nil, err, sout)
	}
}

func BenchmarkExecutePyUnitTests(b *testing.B) {
	p := SubmissionsProcessor{}
	className := "Solution"
	file := className + ".py"
	folder := "testData/adder/"
	unitTest := "\tdef test_is_odd_on_odd(self):\n\t\tself.assertTrue(is_odd(5))"

	pyunitShell := `
import unittest
from %v import *

class %vTestCase(unittest.TestCase):
%v

if __name__ == '__main__':
	unittest.main()`

	pyunitFinal := fmt.Sprintf(string(pyunitShell), className, className, unitTest)

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		p.ExecutePyUnitTests(file, className, folder, pyunitFinal)
	}
}

func TestExecuteJavaStyle(t *testing.T) {
	p := SubmissionsProcessor{}
	folder := "testData/style/"
	file := "Solution.java"
	sout, _, err := p.ExecuteJavaStyle(file, folder)
	if err != nil || strings.Contains(sout, "Audit Done") {
		t.Fatalf("Err should be %v, got %v\n StdOut should contain \"Audit Done\" got %v", nil, err, sout)
	}
}

func BenchmarkExecuteJavaStyle(b *testing.B) {
	p := SubmissionsProcessor{}
	folder := "testData/style/"
	file := "Solution.java"

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		p.ExecuteJavaStyle(file, folder)
	}
}
