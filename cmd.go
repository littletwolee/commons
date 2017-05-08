package commons

import (
	"bytes"
	"fmt"
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
//                  command                *exec.Cmd        cmd point
//                  paras                  []string         parparameters
//                  cmddir                 string           dir of exec
// @Returns output:string err:error
func (c *Cmd) ExecCommand(command string, pars []string, cmddir string) (string, error) {
	cmd := exec.Command(command, pars...)
	if cmddir != "" {
		cmd.Dir = cmddir
	}
	var buf bytes.Buffer
	cmd.Stdout = &buf
	cmd.Stderr = &buf
	err := cmd.Run()
	out := buf.Bytes()
	if err != nil {
		GetLogger().LogErr(err)
	}
	return fmt.Sprintf("%s", out), nil
}
