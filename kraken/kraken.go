package kraken

import (
	"bytes"
	"os/exec"

	"github.com/golang/glog"
)

const (
	K2CLICommand             = "k2cli"
	SubCmdCluster            = "cluster"
	ClusterArgUp             = "up"
	ClusterArgDown           = "down"
	ClusterArgUpdate         = "update"
	UpdateArgAddNodePools    = "--add-nodepools"
	UpdateArgUpdateNodePools = "--update-nodepools"
	UpdateArgRemoveNodePools = "--rm-nodepools"
)

// Execute the command "k2cli" with the specified arguments and returns either;
// on success: the resultant byte array containing stdout, error = nil
// on failure: the resultant byte array containing stderr, error is set
func Execute(command string, commandString []string) ([]byte, error) {
	cmd := exec.Command(command, commandString...)
	stdoutBuf := &bytes.Buffer{}
	stderrBuf := &bytes.Buffer{}
	cmd.Stdout = stdoutBuf
	cmd.Stderr = stderrBuf

	if err := cmd.Run(); err != nil {
		glog.V(3).Infof("k2cli: error return: %v", err)
		return stderrBuf.Bytes(), err
	}
	return stdoutBuf.Bytes(), nil
}
