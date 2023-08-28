package interop

import (
	"bufio"
	"fmt"
	"os/exec"
	"strings"
	"util"
)

func execInternal(name string, arg ...string) (string, error) {
	cmd := exec.Command(name, arg...)
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return "", util.DErr(util.Unavailable, err.Error())
	}
	if err := cmd.Start(); err != nil {
		return "", util.DErr(util.Unavailable, err.Error())
	}
	scanner := bufio.NewScanner(stdout)
	var outputBuilder strings.Builder
	for scanner.Scan() {
		outputBuilder.WriteString(fmt.Sprintf("%v\n", scanner.Text()))
	}
	if err := scanner.Err(); err != nil {
		return "", util.DErr(util.Unavailable, err.Error())
	}
	if err := cmd.Wait(); err != nil {
		return "", util.DErr(util.Unavailable, err.Error())
	}
	return outputBuilder.String(), nil
}
