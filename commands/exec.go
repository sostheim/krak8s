/*
Copyright 2017 Samsung SDSA CNCT

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package commands

import (
	"bytes"
	"os/exec"

	"github.com/golang/glog"
)

// Execute the "command" with the specified arguments and return either;
// on success: the resultant byte array containing stdout, error = nil
// on failure: the resultant byte array containing stderr, error is set
func Execute(command string, commandString []string) ([]byte, error) {
	cmd := exec.Command(command, commandString...)
	stdoutBuf := &bytes.Buffer{}
	stderrBuf := &bytes.Buffer{}
	cmd.Stdout = stdoutBuf
	cmd.Stderr = stderrBuf

	glog.Infof("run cmd:  %s, args: %s", command, commandString)

	if err := cmd.Run(); err != nil {
		glog.Warningf("cmd:  %s, args: %s returned error: %v", command, commandString, err)
		glog.Warningf("cmd:  %s, stderr: %s", command, string(stderrBuf.Bytes()))
		glog.Warningf("cmd:  %s, stdout: %v", command, string(stdoutBuf.Bytes()))
		return stderrBuf.Bytes(), err
	}
	return stdoutBuf.Bytes(), nil
}
