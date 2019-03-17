package processors

import (
	"bytes"
	"log"
	"os/exec"
)

type Processor struct{}

func (jp Processor) Execute(file string, folder string) string {
	var cmd *exec.Cmd
	if folder != "" {
		cmd = exec.Command("dexec", "-C", folder, file)
	} else {
		cmd = exec.Command("dexec", file)
	}
	// cmd.Stdin = strings.NewReader("some input")
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		log.Fatal(err)
	}
	return out.String()
}
