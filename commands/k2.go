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

import "os"

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
	// K2Config - configuration file path
	K2Config = "--config"
	// K2Output - configuration tree output path
	K2Output = "--output"
	// K2AddNodePools - add node pool(s)
	K2AddNodePools = "--addnodepools"
	// K2UpdateNodePools - udpate node pool(s)
	K2UpdateNodePools = "--nodepools"
	// K2RemoveNodePools - remove node pool(s)
	K2RemoveNodePools = "--rmnodepools"
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
	K2ENVK2OptsDefault = "--volume=${KRAKEN}:${KRAKEN} --volume=${SSH_ROOT}:${SSH_ROOT} --volume=${AWS_ROOT}:${AWS_ROOT} -e HOME=${HOME}"
	// K2ENVK2OptsExpansion - K2OPTS environment variable expansion
	K2ENVK2OptsExpansion = "${K2OPTS}"
	// K2ENVExtraVars - KRAKEN_EXTRA_VARS environment variable
	K2ENVExtraVars = "KRAKEN_EXTRA_VARS"
	// K2ExtraVarsAction - extravar
	K2ExtraVarsAction = "kraken_action=update"
	// K2ExtraVarsConfigBase - extravar
	K2ExtraVarsConfigBase = "config_base"
	// K2ExtraVarsConfigPath - extravar
	K2ExtraVarsConfigPath = "config_path"
	// K2ExtraVarsAddNodePools - add node pool(s)
	K2ExtraVarsAddNodePools = "add_nodepools"
	// K2ExtraVarsUpdateNodePools - udpate node pool(s)
	K2ExtraVarsUpdateNodePools = "update_nodepools"
	// K2ExtraVarsRemoveNodePools - remove node pool(s)
	K2ExtraVarsRemoveNodePools = "remove_nodepools"
	// K2VolKraken - KRAKEN volume mount for docker
	K2VolKraken = "--volume=${KRAKEN}:${KRAKEN}"
	// K2VolSSHRoot - SSH_ROOT volume mount for docker
	K2VolSSHRoot = "--volume=${SSH_ROOT}:${SSH_ROOT}"
	// K2VolAWSRoot - AWs_ROOT volume mount for docker
	K2VolAWSRoot = "--volume=${AWS_ROOT}:${AWS_ROOT}"
	// K2EnvHome - HOME environment variable for docker
	K2EnvHome = "-e HOME=${HOME}"
)

// K2SetupEnv - for environment variables that aren't presnet, but required, set them.
func K2SetupEnv() {
	if env := os.Getenv(K2ENVKraken); len(env) == 0 {
		os.Setenv(K2ENVKraken, K2ENVKrakenDefault)
	}
	if env := os.Getenv(K2ENVSSHRoot); len(env) == 0 {
		os.Setenv(K2ENVSSHRoot, K2ENVSSHRootDefault)
	}
	if env := os.Getenv(K2ENVSSHKey); len(env) == 0 {
		os.Setenv(K2ENVSSHKey, K2ENVSSHKeyDefault)
	}
	if env := os.Getenv(K2ENVSSHPub); len(env) == 0 {
		os.Setenv(K2ENVSSHPub, K2ENVSSHPubDefault)
	}
	if env := os.Getenv(K2ENVAWSRoot); len(env) == 0 {
		os.Setenv(K2ENVAWSRoot, K2ENVAWSRootDefault)
	}
	if env := os.Getenv(K2ENVAWSConfig); len(env) == 0 {
		os.Setenv(K2ENVAWSConfig, K2ENVAWSConfigDefault)
	}
	if env := os.Getenv(K2ENVAWSCredentials); len(env) == 0 {
		os.Setenv(K2ENVAWSCredentials, K2ENVAWSCredentialsDefault)
	}
	return
}

// K2EnvString - print environment variables
func K2EnvString() string {
	return K2ENVKraken + "=" + os.Getenv(K2ENVKraken) + ", " +
		K2ENVSSHRoot + "=" + os.Getenv(K2ENVSSHRoot) + ", " +
		K2ENVSSHPub + "=" + os.Getenv(K2ENVSSHPub) + ", " +
		K2ENVSSHKey + "=" + os.Getenv(K2ENVSSHKey) + ", " +
		K2ENVAWSRoot + "=" + os.Getenv(K2ENVAWSRoot) + ", " +
		K2ENVAWSConfig + "=" + os.Getenv(K2ENVAWSConfig) + ", " +
		K2ENVAWSCredentials + "=" + os.Getenv(K2ENVAWSCredentials)
}

// K2EnvExport - environment variables to export (-e) in docker
func K2EnvExport() []string {
	// TODO: reduce this list to a minimal set of required vars
	return []string{"-e " + K2ENVKraken + "=" + os.Getenv(K2ENVKraken),
		"-e " + K2ENVSSHRoot + "=" + os.Getenv(K2ENVSSHRoot),
		"-e " + K2ENVSSHPub + "=" + os.Getenv(K2ENVSSHPub),
		"-e " + K2ENVSSHKey + "=" + os.Getenv(K2ENVSSHKey),
		"-e " + K2ENVAWSRoot + "=" + os.Getenv(K2ENVAWSRoot),
		"-e " + K2ENVAWSConfig + "=" + os.Getenv(K2ENVAWSConfig),
		"-e " + K2ENVAWSCredentials + "=" + os.Getenv(K2ENVAWSCredentials),
	}
}

// K2CmdUp - build a command string to call "./bin/upsh"
func K2CmdUp(docker bool, config string) []string {
	if docker {
		cmd := DockerRunK2()
		cmd = append(cmd, K2Up, config)
		return cmd
	}
	return []string{
		K2Down, config,
	}
}

// K2CmdUpdate - build a command string to call "./bin/update.sh"
func K2CmdUpdate(docker bool, action, base, config, name string) []string {

	nodePoolName := name + "Nodes"
	extraVars := K2ExtraVarsConfigPath + "=" + config + " " +
		K2ExtraVarsConfigBase + "=" + base + " " +
		K2ExtraVarsAction + " " + action + "=" + nodePoolName
	os.Setenv(K2ENVExtraVars, extraVars)

	var k2UpdateArg string
	if action == K2ExtraVarsAddNodePools {
		k2UpdateArg = K2AddNodePools
	} else if action == K2ExtraVarsUpdateNodePools {
		k2UpdateArg = K2UpdateNodePools
	} else if action == K2ExtraVarsRemoveNodePools {
		k2UpdateArg = K2RemoveNodePools
	} // TODO: else { blow up! }

	if docker {
		cmd := DockerRunK2()
		cmd = append(cmd, K2Update, K2Config, config, K2Output, base, k2UpdateArg, nodePoolName)
		return cmd
	}
	return []string{
		K2Update,
	}
}

// K2CmdDown - build a command string to call "./bin/down.sh"
func K2CmdDown(docker bool, config string) []string {
	if docker {
		cmd := DockerRunK2()
		cmd = append(cmd, K2Down, config)
		return cmd
	}
	return []string{
		K2Down, config,
	}
}
