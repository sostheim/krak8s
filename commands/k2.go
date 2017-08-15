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
	// K2 - kraken command name
	K2 = "k2"
	// K2Image - kraken container image
	K2Image = "quay.io/samsung_cnct/k2:latest"
	// K2Up -  "up"
	K2Up = "./bin/up.sh"
	// K2Down -  "down"
	K2Down = "./bin/down.sh"
	// K2Update -  "update"
	K2Update = "./bin/update.sh"
	// K2ENVKraken - KRAKEN environment variable
	K2ENVKraken = "KRAKEN"
	// K2ENVKrakenDefault - KRAKEN environment variable default value
	K2ENVKrakenDefault = "${HOME}/.kraken"
	// K2ENVSSHRoot - SSH_ROOT environment variable
	K2ENVSSHRoot = "SSH_ROOT"
	// K2ENVSSHRootDefault - SSH_ROOT environment variable default value
	K2ENVSSHRootDefault = "${HOME}/.ssh"
	// K2ENVAWSRoot - AWS_ROOT environment variable
	K2ENVAWSRoot = "AWS_ROOT"
	// K2ENVAWSRootDefault - AWS_ROOT environment variable default value
	K2ENVAWSRootDefault = "${HOME}/.aws"
	// K2ENVAWSConfig - AWS_CONFIG environment variable
	K2ENVAWSConfig = "AWS_CONFIG"
	// K2ENVAWSConfigDefault - AWS_CONFIG environment variable default value
	K2ENVAWSConfigDefault = "${AWS_ROOT}/config"
	// K2ENVAWSCredentials - AWS_CREDENTIALS environment variable
	K2ENVAWSCredentials = "AWS_CREDENTIALS"
	// K2ENVAWSCredentialsDefault - AWS_CREDENTIALS environment variable default value
	K2ENVAWSCredentialsDefault = "${AWS_ROOT}/credentials"
	// K2ENVSSHKey - SSH_KEY environment variable
	K2ENVSSHKey = "SSH_KEY"
	// K2ENVSSHKeyDefault - SSH_KEY environment variable default value
	K2ENVSSHKeyDefault = "${SSH_ROOT}/id_rsa"
	// K2ENVSSHPub - SSH_PUB environment variable
	K2ENVSSHPub = "SSH_PUB"
	// K2ENVSSHPubDefault - SSH_PUB environment variable default value
	K2ENVSSHPubDefault = "${SSH_ROOT}/id_rsa.pub"
	// K2ENVK2Opts - K2OPTS environment variable
	K2ENVK2Opts = "K2OPTS"
	// K2ENVK2OptsDefault - KRAKEN environment variable default value
	K2ENVK2OptsDefault = "-v ${KRAKEN}:${KRAKEN} -v ${SSH_ROOT}:${SSH_ROOT} -v ${AWS_ROOT}:${AWS_ROOT} -e HOME=${HOME} --rm=true -it"
	// K2ENVK2OptsExpansion - K2OPTS environment variable expansion
	K2ENVK2OptsExpansion = "${K2OPTS}"
)

// Up - build a command string to call "./bin/upsh"
func Up(docker bool, config string) []string {
	if docker {
		cmd := DockerRunCmd()
		cmd = append(cmd, K2ENVK2OptsDefault, K2Image, K2Up, config)
		return cmd
	}
	return []string{
		K2Update, config,
	}
}

// Update - build a command string to call "./bin/update.sh"
func Update(docker bool, config string) []string {
	if docker {
		cmd := DockerRunCmd()
		cmd = append(cmd, K2ENVK2OptsDefault, K2Image, K2Update, config)
		return cmd
	}
	return []string{
		K2Update, config,
	}
}

// Down - build a command string to call "./bin/down.sh"
func Down(docker bool, config string) []string {
	if docker {
		cmd := DockerRunCmd()
		cmd = append(cmd, K2ENVK2OptsExpansion, K2Image, K2Down, config)
		return cmd
	}
	return []string{
		K2Update, config,
	}
}
