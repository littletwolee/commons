package commons

import (
	"bytes"
	"errors"
	"fmt"
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
//                  command                *exec.Cmd        cmd point
//                  paras                  []string         parparameters
//                  cmddir                 string           dir of exec
// @Returns output:string err:error
func (c *Cmd) ExecCommand(command string, pars []string, cmddir string) (string, error) {
	cmd := exec.Command(command, pars...)
	if cmddir != "" {
		cmd.Dir = cmddir
	}
	var resultbuf bytes.Buffer
	var errbuf bytes.Buffer
	cmd.Stdout = &resultbuf
	cmd.Stderr = &errbuf
	err := cmd.Run()
	outresultbuf := resultbuf.Bytes()
	outerrbuf := errbuf.Bytes()
	if err != nil {
		return "", err
	}
	if len(outerrbuf) > 0 {
		err = errors.New(fmt.Sprintf("%s", outerrbuf))
		return "", err
	}
	return fmt.Sprintf("%s", outresultbuf), nil
}
