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
	"fmt"
	"krak8s/commands"
	"os"
	"strings"

	"github.com/golang/glog"
	flag "github.com/spf13/pflag"
)

type config struct {
	flagSet          *flag.FlagSet
	kubeconfig       *string
	proxy            *string
	serviceName      *string
	version          *bool
	healthCheck      *bool
	krakenConfigFile *string
	krakenConfigDir  *string
	krakenKeyPair    *string
	krakenKubeConfig *string
	krakenCommand    *string
	dryrun           *bool
	debug            *bool
}

func newConfig() *config {
	return &config{
		kubeconfig:       flag.String("kubeconfig", "", "absolute path to the kubeconfig file"),
		proxy:            flag.String("proxy", "", "kubctl proxy server running at the given url"),
		version:          flag.Bool("version", false, "display version info and exit"),
		healthCheck:      flag.Bool("health-check", false, "enable health checking for API service"),
		krakenConfigFile: flag.String("kraken-config-file", commands.DefaultConfigFile, "kraken configuration yaml file name"),
		krakenConfigDir:  flag.String("kraken-config-dir", commands.DefaultConfigDir, "kraken configuration yaml directory path"),
		krakenKeyPair:    flag.String("kraken-nodepool-keypair", commands.DefaultKeyPair, "kraken configuration yaml: deployment.clusters[0].nodePools.keyPair"),
		krakenKubeConfig: flag.String("kraken-kubeconfig", commands.DefaultKubeConfig, "kraken confiuration yaml: deployment.clusters[0].nodePools.kubeConfig"),
		krakenCommand:    flag.String("kraken-command", commands.K2, "command to run to execute kraken operations, either `k2`, or `k2cli` only"),
		dryrun:           flag.Bool("dry-run", false, "don't actually execute backend commands"),
		debug:            flag.Bool("debug", false, "enable debug output"),
	}
}

func (cfg *config) String() string {
	return fmt.Sprintf("Configuration Data: kubeconfig: %s, proxy: %s, "+
		"health-check: %t, version: %t, kraken-config-file: %s, "+
		"kraken-config-dir: %s, kraken-nodepool-keypair: %s, "+
		"kraken-kubeconfig: %s, kraken-command: %s, dry-run: %t",
		*cfg.kubeconfig, *cfg.proxy, *cfg.healthCheck, *cfg.version,
		*cfg.krakenConfigFile, *cfg.krakenConfigDir, *cfg.krakenKeyPair,
		*cfg.krakenKubeConfig, *cfg.krakenCommand, *cfg.dryrun)
}

// For any configuration members that contain environment variables as values, expand them.
func (cfg *config) envExpand() {
	*cfg.kubeconfig = os.ExpandEnv(*cfg.kubeconfig)
	*cfg.krakenConfigFile = os.ExpandEnv(*cfg.krakenConfigFile)
	*cfg.krakenConfigDir = os.ExpandEnv(*cfg.krakenConfigDir)
	*cfg.krakenKeyPair = os.ExpandEnv(*cfg.krakenKeyPair)
	*cfg.krakenKubeConfig = os.ExpandEnv(*cfg.krakenKubeConfig)
	*cfg.krakenCommand = os.ExpandEnv(*cfg.krakenCommand)
}

var envSupport = map[string]bool{
	"kubeconfig":              true,
	"proxy":                   true,
	"version":                 false,
	"health-check":            true,
	"kraken-config-file":      true,
	"kraken-config-dir":       true,
	"kraken-nodepool-keypair": true,
	"kraken-kubeconfig":       true,
	"kraken-command":          true,
	"dry-run":                 false,
	"debug":                   false,
}

func variableName(name string) string {
	return "KRAK8S_" + strings.ToUpper(strings.Replace(name, "-", "_", -1))
}

// Just like Flags.Parse() except we try to get recognized values for the valid
// set of flags from environment variables.  We choose to use the environment
// value if 1) the value hasen't already been set as command line flags and the
// flas is a member of the supported set (see map defined above).
func (cfg *config) envParse() error {
	var err error

	alreadySet := make(map[string]bool)
	cfg.flagSet.Visit(func(f *flag.Flag) {
		if envSupport[f.Name] {
			alreadySet[variableName(f.Name)] = true
		}
	})

	usedEnvKey := make(map[string]bool)
	cfg.flagSet.VisitAll(func(f *flag.Flag) {
		if envSupport[f.Name] {
			key := variableName(f.Name)
			if !alreadySet[key] {
				val := os.Getenv(key)
				if val != "" {
					usedEnvKey[key] = true
					if serr := cfg.flagSet.Set(f.Name, val); serr != nil {
						err = fmt.Errorf("invalid value %q for %s: %v", val, key, serr)
					}
					glog.V(3).Infof("recognized and used environment variable %s=%s", key, val)
				}
			}
		}
	})

	for _, env := range os.Environ() {
		kv := strings.SplitN(env, "=", 2)
		if len(kv) != 2 {
			glog.Warningf("found invalid env %s", env)
		}
		if usedEnvKey[kv[0]] {
			continue
		}
		if alreadySet[kv[0]] {
			glog.V(3).Infof("recognized environment variable %s, but unused: superseeded by command line flag ", kv[0])
			continue
		}
		if strings.HasPrefix(env, "KRAK8S_") {
			glog.Warningf("unrecognized environment variable %s", env)
		}
	}

	return err
}
