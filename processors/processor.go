package processors

import (
	"bytes"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"strings"
)

type Processor struct{}

func (p Processor) Execute(file string, folder string) string {
	var cmd *exec.Cmd
	if folder != "" {
		cmd = exec.Command("dexec", "-C", folder, file)
	} else {
		cmd = exec.Command("dexec", file)
	}
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		log.Fatal(err)
	}
	return out.String()
}

func (p Processor) ExecuteWithInput(file string, folder string, input string) string {
	// compile command
	var cmd *exec.Cmd
	if folder != "" {
		cmd = exec.Command("dexec", "-C", folder, file)
	} else {
		cmd = exec.Command("dexec", file)
	}
	// provide stdin
	cmd.Stdin = strings.NewReader(input)
	// take back stdout
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		log.Fatal(err)
	}
	return out.String()
}

func (p Processor) ExecuteJUnitTests(className string, folder string, junitTests string) string {
	fileName := className + "Test.java"
	path := folder + fileName
	// detect if file exists
	_, err := os.Stat(path)
	// delete file if exists
	if os.IsExist(err) {
		err = os.Remove(path)
		if err != nil {
			panic(err)
		}
	}
	err = ioutil.WriteFile(path, []byte(junitTests), 0644)
	if err != nil {
		panic(err)
	}

	copyJunit := exec.Command("cp", "assets/junit.jar", folder)
	err = copyJunit.Run()
	if err != nil {
		panic(err)
	}
	compileClass := exec.Command("javac", className+".java")
	compileClass.Dir = folder
	err = compileClass.Run()
	if err != nil {
		panic(err)
	}
	compileTestClass := exec.Command("javac", "-cp", ".:junit.jar", fileName)
	compileTestClass.Dir = folder
	err = compileTestClass.Run()
	if err != nil {
		panic(err)
	}
	runTestClass := exec.Command("java", "-cp", ".:junit.jar", "org.junit.runner.JUnitCore", className+"Test")
	runTestClass.Dir = folder
	var out bytes.Buffer
	runTestClass.Stdout = &out
	if err = runTestClass.Run(); err != nil {
		panic(err)
	}
	return out.String()
}
