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
	// KCTL - kubectl command name
	KCTL = "KCTL"
	// KCTLGet - "get"
	KCTLGet = "get"
	// KCTLNodes - "nodes"
	KCTLNodes = "nodes"
	// KCTLPods - "pods"
	KCTLPods = "pods"
	// KCTLNS - "ns" namespaces
	KCTLNS = "ns"
	// KCTLArgF - filename, directory, or URL to files
	KCTLArgF = "-f"
	// KCTLArgL - label selector
	KCTLArgL = "-l"
	// KCTLArgN - argment for a namespaces
	KCTLArgN = "-n"
	// KCTLArgAll - argment for all output from command
	KCTLArgAll = "--all"
	// KCTLArgAllNS - argment for all-namespaces
	KCTLArgAllNS = "--all-namespaces"
	// KCTLArgJSON - argment for JSON command output
	KCTLArgJSON = "-ojson"
)
