// Package test contains all the test tasks exposed by the grader
package test

// Test is a test task containing all the required properties
type Test struct {
	Type           string   `json:"type"`           // required - "io" / "output" / "junit" / "pyunit" / "javaStyle"
	Description    string   `json:"description"`    // optional - Any description of your task
	Input          []string `json:"input"`          // Used by "io" - List of strings where they'd be joined by a new line character when provided as a stdin buffer
	ExpectedOutput []string `json:"expectedOutput"` // Used by "output" and "io" - List of strings representing the expected output after being split by the new line char
	UnitTest       string   `json:"unitTest"`
}

// Result represents all the properties returned from runnning
// a test task
type Result struct {
	Test        Test    `json:"test"`        // The test executed
	Description string  `json:"description"` // Description of the test result (if any)
	StdOut      string  `json:"stdOut"`      // Standard output
	StdErr      string  `json:"stdErr"`      // Standard error
	Successful  bool    `json:"successful"`  // Was it a successful test execution?
	Similarity  float64 `json:"similarity"`  // If applicable, how similar was the standard ouptut to the expected output (jaro winkler)
	TimedOut    bool    `json:"timedOut"`    // Did the test execution time out
}

// Task interface defines methods that need
// to be implemented for a test task handler to be considered
// of type TestTask
type Task interface {
	RunTest() (Result, error)
}
