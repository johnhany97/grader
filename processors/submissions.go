package processors

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"
)

type SubmissionsProcessor struct{}

func (p SubmissionsProcessor) Execute(file string, folder string) (string, string) {
	var cmd *exec.Cmd
	if folder != "" {
		cmd = exec.Command("dexec", "-C", folder, file)
	} else {
		cmd = exec.Command("dexec", file)
	}
	var out bytes.Buffer
	var e bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &e
	cmd.Run()
	return out.String(), e.String()
}

func (p SubmissionsProcessor) ExecuteWithInput(file string, folder string, input string) (string, string) {
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
	var e bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &e
	cmd.Run()
	return out.String(), e.String()
}

func (p SubmissionsProcessor) ExecuteJUnitTests(className string, folder string, junitTests string) (string, error) {
	fileName := className + "Test.java"
	path := folder + fileName
	// detect if file exists
	_, err := os.Stat(path)
	// delete file if exists
	if os.IsExist(err) {
		err = os.Remove(path)
		if err != nil {
			return "", err
		}
	}
	err = ioutil.WriteFile(path, []byte(junitTests), 0644)
	if err != nil {
		return "", err
	}
	// delete when done
	defer func() {
		var err = os.Remove(path)
		if err != nil {
			return
		}
	}()
	junitCmd := exec.Command("/bin/sh", "-c", fmt.Sprintf("docker run -t --rm -v $(pwd -P)/%v.java:/tmp/dexec/build/%v.java -v $(pwd -P)/%v:/tmp/dexec/build/%v johnhany97/grader-junit %v.java %v", className, className, fileName, fileName, className, fileName))
	junitCmd.Dir = folder
	var out bytes.Buffer
	junitCmd.Stdout = &out
	if err = junitCmd.Run(); err != nil {
		return "", err
	}
	return out.String(), nil
}

func (p SubmissionsProcessor) ExecutePyUnitTests(file string, className string, folder string, pyUnitTests string) (string, error) {
	fileName := "test_" + file
	path := folder + fileName
	// detect if file exists
	_, err := os.Stat(path)
	// delete file if exists
	if os.IsExist(err) {
		err = os.Remove(path)
		if err != nil {
			return "", err
		}
	}
	err = ioutil.WriteFile(path, []byte(pyUnitTests), 0644)
	if err != nil {
		return "", err
	}
	// delete when done
	defer func() {
		var err = os.Remove(path)
		if err != nil {
			return
		}
	}()
	var cmd *exec.Cmd
	cmd = exec.Command("dexec", fileName, "-i", file)
	cmd.Dir = folder
	var out bytes.Buffer
	cmd.Stderr = &out
	if err = cmd.Run(); err != nil {
		return "", err
	}
	return out.String(), nil
}
