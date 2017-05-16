package commons

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
)

var (
	consCmd *Cmd
)

type Cmd struct{}

func GetCmd() *Cmd {
	if consCmd == nil {
		consCmd = &Cmd{}
	}
	return consCmd
}

// @Title ExecCommand
// @Description exec command
// @Parameters
//                  commandName            *exec.Cmd        cmd point
//                  rootdir                string           dir of exec
//                  paras                  []string         parparameters
//                  isOutput               bool             is output in os
// @Returns output:string outerr:string err:error
func (c *Cmd) ExecCommand(commandName, rootDir string, params []string, isOutput bool) (string, string, error) {
	var (
		cmd *exec.Cmd
	)
	if params == nil || len(params) == 0 {
		cmd = exec.Command(commandName)
	} else {
		cmd = exec.Command(commandName, params...)
	}
	cmd.Dir = rootDir
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return "", "", err
	}
	stderr, err := cmd.StderrPipe()
	if err != nil {
		return "", "", err
	}
	if isOutput {
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
	}
	cmd.Start()
	readerout := bufio.NewReader(stdout)
	strout, err := ioutil.ReadAll(readerout)
	if err != nil {
		return "", "", err
	}
	readererr := bufio.NewReader(stderr)
	strerr, err := ioutil.ReadAll(readererr)
	if err != nil {
		return "", "", err
	}
	cmd.Wait()
	return fmt.Sprintf("%s", strout), fmt.Sprintf("%s", strerr), err
}
