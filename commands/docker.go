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

const (
	// Docker - docker cli command name
	Docker = "docker"
	// DockerRun - docker command "run"
	DockerRun = "run"
	// DockerPull - docker command "pull"
	DockerPull = "pull"
	// DockerPush - docker command "push"
	DockerPush = "push"
	// DockerPS - docker command "ps"
	DockerPS = "ps"
	// DockerImages - docker command "images"
	DockerImages = "images"
	// DockerRemoveImages - docker command "images"
	DockerRemoveImages = "rmi"
	// DockerArgAll - argment for all output from command
	DockerArgAll = "--all"
	// DockerArgForce - argment to force an action
	DockerArgForce = "--force"
	// DockerArgIT - argment for console (i - Keep STDIN open even if not attached) (t - Allocate a pseudo-tty)
	DockerArgIT = "-it"
	// DockerArgRM - Automatically remove the container when it exits
	DockerArgRM = "--rm"
	// DockerArgFormatJSON - generate JSON output from command
	DockerArgFormatJSON = `--format "{{json .}}"`
)

// DockerRunCmd - return the common docker run command
func DockerRunCmd() []string {
	return []string{
		Docker, DockerRun,
	}
}

// DockerRunK2 - docker run command with args set for K2
// common prefix: docker run $K2OPTS quay.io/samsung_cnct/k2:latest
func DockerRunK2() []string {
	run := []string{
		Docker, DockerRun, DockerArgRM, K2VolKraken, K2VolAWSRoot, K2VolSSHRoot, K2EnvHome,
	}
	run = append(run, K2DockerEnvExport()...)
	run = append(run, K2Image)
	return run
}

// DockerRunIT - docker run command with -it args
func DockerRunIT() []string {
	return []string{
		Docker, DockerRun, DockerArgIT,
	}
}

// DockerRunITRM - docker run command with -it -rm args
func DockerRunITRM() []string {
	return []string{
		Docker, DockerRun, DockerArgRM,
	}
}

// DockerRunRM - docker run command with -rm args
func DockerRunRM() []string {
	return []string{
		Docker, DockerRun, DockerArgRM,
	}
}
