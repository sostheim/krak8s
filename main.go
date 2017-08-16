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

package main

import (
	goflag "flag"
	"fmt"
	"krak8s/commands"
	"time"

	"github.com/blang/semver"
	"github.com/golang/glog"
	flag "github.com/spf13/pflag"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/pkg/api/unversioned"
	"k8s.io/client-go/pkg/util/wait"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

// MajorMinorPatch - semantic version string
var MajorMinorPatch string

// ReleaseType - release type
var ReleaseType = "alpha"

// GitCommit - git commit sha-1 hash
var GitCommit string

var krak8sCfg *config

func init() {
	go wait.Until(glog.Flush, 10*time.Second, wait.NeverStop)
	krak8sCfg = newConfig()
}

func addGV(config *rest.Config) {
	config.ContentConfig.GroupVersion = &unversioned.GroupVersion{
		Group:   "",
		Version: "v1",
	}
}

func inCluster() *rest.Config {
	glog.V(3).Infof("inCluster(): creating config")
	config, err := rest.InClusterConfig()
	if err != nil {
		panic(err.Error())
	}
	return config
}

func external() *rest.Config {
	glog.V(3).Infof("external(): creating config")
	config, err := clientcmd.BuildConfigFromFlags("", *krak8sCfg.kubeconfig)
	if err != nil {
		panic(err.Error())
	}
	return config
}

func byProxy() *rest.Config {
	glog.V(3).Infof("byProxy(): creating config")
	return &rest.Config{
		Host: *krak8sCfg.proxy,
	}
}

func displayVersion() {
	semVer, err := semver.Make(MajorMinorPatch + "-" + ReleaseType + "+git.sha." + GitCommit)
	if err != nil {
		panic(err)
	}
	fmt.Println(semVer.String())
}

func main() {
	flag.CommandLine.AddGoFlagSet(goflag.CommandLine)
	flag.Parse()

	krak8sCfg.flagSet = flag.CommandLine

	// check for version flag, if present print veriosn and exit
	if *krak8sCfg.version {
		displayVersion()
		return
	}

	glog.Infof("main(): initial configuration: %v", krak8sCfg.String())
	krak8sCfg.envParse()
	glog.Infof("main(): env override configuration: %v", krak8sCfg.String())
	krak8sCfg.envExpand()
	glog.Infof("main(): env expanded configuration: %v", krak8sCfg.String())

	RunnerSetup()
	commands.SetDebug(*krak8sCfg.debug)
	commands.SetDryrun(*krak8sCfg.dryrun)

	// creates the config, in preference order, for:
	// 1 - the proxy URL, if present as an argument
	// 2 - kubeconfig, if present as an argument
	// 3 - otherwise assume execution on an in-cluster node
	//     note: this will fail with the appropriate error messages
	//           if not actually executing on a node in the cluster.
	var config *rest.Config
	if *krak8sCfg.proxy != "" {
		config = byProxy()
	} else if *krak8sCfg.kubeconfig != "" {
		config = external()
	} else {
		config = inCluster()
	}
	addGV(config)

	// creates a clientset
	glog.V(3).Infof("main(): create clientset from config")
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}

	backend := NewRunner()
	go backend.ProcessRequests()

	// Create our REST Service's API Server
	glog.V(3).Infof("main(): staring API service")
	srv := newAPIServer(clientset, krak8sCfg, backend)

	// Start service server
	srv.run()
}
