package cmd

import (
	"buding-job/common/log"
	"github.com/spf13/cobra"
	"os/exec"
)

func ExecuteCommand(name string, subname string, args ...string) (string, error) {
	args = append([]string{subname}, args...)

	cmd := exec.Command(name, args...)
	bytes, err := cmd.CombinedOutput()

	return string(bytes), err
}

func Error(cmd *cobra.Command, args []string, err error) {
	log.GetLog().Fatalf("execute %s args:%v error:%v\n", cmd.Name(), args, err)
}
