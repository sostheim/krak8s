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
	// K2CLI - kraken cli command name
	K2CLI = "k2cli"
	// K2CLICluster - k2cli subcommand
	K2CLICluster = "cluster"
	// K2CLIClusterUp - cluster subcommand argment "up"
	K2CLIClusterUp = "up"
	// K2CLIClusterDown - cluster subcommand argment "down"
	K2CLIClusterDown = "down"
	// K2CLIClusterUpdate - cluster subcommand argment "update"
	K2CLIClusterUpdate = "update"
	// K2CLIAddNodePools - cluster subcommand update argment to add node pool
	K2CLIAddNodePools = "--add-nodepools"
	// K2CLIUpdateNodePools - cluster subcommand update argment to update node pool
	K2CLIUpdateNodePools = "--update-nodepools"
	// K2CLIRemoveNodePools - cluster subcommand update argment to remove node pool
	K2CLIRemoveNodePools = "--rm-nodepools"
)

// ClusterUpdateAdd - build a command string to call "cluster update --add-nodepools"
func ClusterUpdateAdd(name string) []string {
	return []string{
		K2CLI, K2CLICluster, K2CLIClusterUpdate, K2CLIAddNodePools, name + "Nodes",
	}
}

// ClusterUpdateRemove - build a command string to call "cluster update --rm-nodepools"
func ClusterUpdateRemove(name string) []string {
	return []string{
		K2CLI, K2CLICluster, K2CLIClusterUpdate, K2CLIRemoveNodePools, name + "Nodes",
	}
}
