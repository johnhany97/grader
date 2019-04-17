package test

type Test struct {
	Type           string   `json:"type"`
	Description    string   `json:"description"`
	Input          []string `json:"input"`
	ExpectedOutput []string `json:"expectedOutput"`
	UnitTest       string   `json:"unitTest"`
}

type TestResult struct {
	Test        Test    `json:"test"`
	Description string  `json:"description"`
	StdOut      string  `json:"stdOut"`
	StdErr      string  `json:"stdErr"`
	Successful  bool    `json:"successful"`
	Similarity  float64 `json:"similarity"`
}

type TestType int

const (
	InputOutput TestType = iota
	Output
	FileInputOutput
	JUnit
)

type TestTask interface {
	RunTest() (TestResult, error)
}
