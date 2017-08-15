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
	"os"
	"os/exec"

	"github.com/golang/glog"
)

// check all command arguments to see if they are environment variables
func envCheck(args []string) []string {

	expanded := make([]string, len(args))
	copy(expanded, args)

	done := false
	i := 4 // maximum depth of expansions incase of circular or recursive definition
	for !done && i > 0 {
		i--
		done = true
		for index, arg := range expanded {
			exp := os.ExpandEnv(arg)
			if exp != arg {
				done = false
				expanded[index] = exp
			}
		}
	}
	return expanded
}

// Execute the "command" with the specified arguments and return either;
// on success: the resultant byte array containing stdout, error = nil
// on failure: the resultant byte array containing stderr, error is set
func Execute(command string, arguments []string, dryrun bool) ([]byte, error) {
	expandedArguments := envCheck(arguments)

	cmd := exec.Command(command, expandedArguments...)
	stdoutBuf := &bytes.Buffer{}
	stderrBuf := &bytes.Buffer{}
	cmd.Stdout = stdoutBuf
	cmd.Stderr = stderrBuf

	glog.Infof("run cmd:  %s, args: %s", command, arguments)
	glog.Infof("run cmd:  %s, env ${args}: %s", command, expandedArguments)

	if dryrun {
		return stdoutBuf.Bytes(), nil
	}

	if err := cmd.Run(); err != nil {
		glog.Warningf("cmd:  %s, args: %s returned error: %v", command, arguments, err)
		glog.Warningf("cmd:  %s, stderr: %s", command, string(stderrBuf.Bytes()))
		glog.Warningf("cmd:  %s, stdout: %v", command, string(stdoutBuf.Bytes()))
		return stderrBuf.Bytes(), err
	}
	return stdoutBuf.Bytes(), nil
}
