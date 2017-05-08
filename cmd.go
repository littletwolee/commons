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
//                  cmd             *exec.Cmd        cmd point
//                  isStdout        bool             need os out put
// @Returns output:string err:error
func (c *Cmd) ExecCommand(cmd *exec.Cmd, isStdout bool) (string, error) {
	var buf bytes.Buffer
	cmd.Stdout = &buf
	cmd.Stderr = &buf
	if isStdout {
		cmd.Stdout = os.Stdout // 重定向标准输出
		cmd.Stderr = os.Stderr // 重定向标准输出
	}
	err := cmd.Run()
	out := buf.Bytes()
	if err != nil {
		GetLogger().LogErr(err)
	}
	return fmt.Sprintf("%s", out), nil
}
