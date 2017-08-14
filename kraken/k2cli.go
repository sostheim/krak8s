package kraken

import (
	"bytes"
	"os/exec"

	"github.com/golang/glog"
)

const (
	// K2CLICommand - kraken cli command name
	K2CLICommand = "k2cli"
	// SubCmdCluster - k2cli subcommand
	SubCmdCluster = "cluster"
	// ClusterArgUp - cluster subcommand argment "up"
	ClusterArgUp = "up"
	// ClusterArgDown - cluster subcommand argment "down"
	ClusterArgDown = "down"
	// ClusterArgUpdate - cluster subcommand argment "update"
	ClusterArgUpdate = "update"
	// UpdateArgAddNodePools - cluster subcommand update argment to add node pool
	UpdateArgAddNodePools = "--add-nodepools"
	// UpdateArgUpdateNodePools - cluster subcommand update argment to update node pool
	UpdateArgUpdateNodePools = "--update-nodepools"
	// UpdateArgRemoveNodePools - cluster subcommand update argment to remove node pool
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

	glog.Infof("k2cli.Execute(): cmd:  %s, args: %s", command, commandString)

	if err := cmd.Run(); err != nil {
		glog.Warningf("k2cli.Execute(): cmd:  %s, args: %s returned error: %v", command, commandString, err)
		glog.Warningf("k2cli.Execute(): cmd:  %s, stderr: %s", command, string(stderrBuf.Bytes()))
		glog.Warningf("k2cli.Execute(): cmd:  %s, stdout: %v", command, string(stdoutBuf.Bytes()))
		return stderrBuf.Bytes(), err
	}
	return stdoutBuf.Bytes(), nil
}
