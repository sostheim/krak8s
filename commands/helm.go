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
	// Helm - helm command name
	Helm = "helm"
	// HelmDelete - Helm subcommand delete
	HelmDelete = "delete"
	// HelmFetch - Helm subcommand fetch
	HelmFetch = "fetch"
	// HelmGet - Helm subcommand get
	HelmGet = "get"
	// HelmInstall - Helm subcommand install
	HelmInstall = "install"
	// HelmStatus - Helm subcommand status
	HelmStatus = "status"
	// HelmUpdate - Helm subcommand update
	HelmUpdate = "update"
	// HelmUpgrade - Helm subcommand upgrade
	HelmUpgrade = "upgrade"
	// HelmArgKubeContext - name of the kubeconfig context to use
	HelmArgKubeContext = "--kube-context"
	// HelmArgName - chart/release name
	HelmArgName = "--name"
	// HelmArgNamespace - k8s namespace name
	HelmArgNamespace = "--namespace"
	// HelmArgSet - set values on the command line
	HelmArgSet = "--set"
	// HelmArgTillerNS - namespace of tiller (default "kube-system")
	HelmArgTillerNS = "--tiller-namespace"
	// HelmArgTimeout - time in seconds to wait for any individual operation (default 300)
	HelmArgTimeout = "--timeout"
	// HelmArgVersion - chart version
	HelmArgVersion = "--version"
	// HelmArgWait - wait until all elements are created
	HelmArgWait = "--wait"
)
