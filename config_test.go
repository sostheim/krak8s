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

	flag "github.com/spf13/pflag"

	"krak8s/commands"
	"testing"
)

func validateBoolFlag(flag string, value bool, b *bool, t *testing.T) bool {
	if b == nil {
		t.Errorf("validateBoolFlag(%s): want object, have nil", flag)
		return false
	} else if *b != value {
		t.Errorf("validateBoolFlag(%s): want %t, have %t", flag, value, *b)
		return false
	}
	return true
}

func validateStringFlag(flag, value string, s *string, t *testing.T) bool {
	if s == nil {
		t.Errorf("validateStringFlag(%s): want object, have nil", flag)
		return false
	} else if *s != value {
		t.Errorf("validateStringFlag(%s): want %s, have %s", flag, value, *s)
		return false
	}
	return true
}

func TestNewConfig(t *testing.T) {
	cfg := krak8sCfg
	f := flag.NewFlagSet("test", flag.ContinueOnError)
	f.AddGoFlagSet(goflag.CommandLine)

	if !validateStringFlag("kubeconfig", "", cfg.kubeconfig, t) {
		t.Error("TestNewConfig() want valid kubeconfig")
	}
	if !validateStringFlag("proxy", "", cfg.proxy, t) {
		t.Error("TestNewConfig() want valid proxy")
	}
	if !validateBoolFlag("version", false, cfg.version, t) {
		t.Error("TestNewConfig() want valid version")
	}
	if !validateBoolFlag("healthCheck", true, cfg.healthCheck, t) {
		t.Error("TestNewConfig() want valid healthCheck")
	}
	if !validateStringFlag("krakenConfigFile", commands.DefaultConfigFile, cfg.krakenConfigFile, t) {
		t.Error("TestNewConfig() want valid krakenConfigFile")
	}
	if !validateStringFlag("krakenConfigDir", commands.DefaultConfigDir, cfg.krakenConfigDir, t) {
		t.Error("TestNewConfig() want valid krakenConfigDir")
	}
	if !validateStringFlag("krakenKeyPair", commands.DefaultKeyPair, cfg.krakenKeyPair, t) {
		t.Error("TestNewConfig() want valid krakenKeyPair")
	}
	if !validateStringFlag("krakenKubeConfig", commands.DefaultKubeConfig, cfg.krakenKubeConfig, t) {
		t.Error("TestNewConfig() want valid krakenKubeConfig")
	}
	if !validateStringFlag("krakenCommand", commands.K2, cfg.krakenCommand, t) {
		t.Error("TestNewConfig() want valid krakenCommand")
	}
	if !validateBoolFlag("debug", false, cfg.debug, t) {
		t.Error("TestNewConfig() want valid debug")
	}
	if !validateBoolFlag("dryrun", false, cfg.dryrun, t) {
		t.Error("TestNewConfig() want valid dryrun")
	}
}

func TestConfigOverrides(t *testing.T) {
	cfg := krak8sCfg
	f := flag.NewFlagSet("test", flag.ContinueOnError)
	f.AddGoFlagSet(goflag.CommandLine)

	flag.Set("kubeconfig", "/Users/sostheim/.kube/config")
	flag.Set("kraken-config-file", "config.yaml")
	flag.Set("kraken-kubeconfig", "geographKube")
	flag.Set("health-check", "true")
	flag.Set("dry-run", "true")
	flag.Parse()

	if !validateStringFlag("kubeconfig", "/Users/sostheim/.kube/config", cfg.kubeconfig, t) {
		t.Error("TestConfigOverrides() want valid kubeconfig")
	}
	if !validateStringFlag("proxy", "", cfg.proxy, t) {
		t.Error("TestConfigOverrides() want valid proxy")
	}
	if !validateBoolFlag("version", false, cfg.version, t) {
		t.Error("TestConfigOverrides() want valid version")
	}
	if !validateBoolFlag("healthCheck", true, cfg.healthCheck, t) {
		t.Error("TestConfigOverrides() want valid healthCheck")
	}
	if !validateStringFlag("krakenConfigFile", "config.yaml", cfg.krakenConfigFile, t) {
		t.Error("TestConfigOverrides() want valid krakenConfigFile")
	}
	if !validateStringFlag("krakenConfigDir", commands.DefaultConfigDir, cfg.krakenConfigDir, t) {
		t.Error("TestConfigOverrides() want valid krakenConfigDir")
	}
	if !validateStringFlag("krakenKeyPair", commands.DefaultKeyPair, cfg.krakenKeyPair, t) {
		t.Error("TestConfigOverrides() want valid krakenKeyPair")
	}
	if !validateStringFlag("krakenKubeConfig", "geographKube", cfg.krakenKubeConfig, t) {
		t.Error("TestConfigOverrides() want valid krakenKubeConfig")
	}
	if !validateStringFlag("krakenCommand", commands.K2, cfg.krakenCommand, t) {
		t.Error("TestConfigOverrides() want valid krakenCommand")
	}
	if !validateBoolFlag("debug", false, cfg.debug, t) {
		t.Error("TestConfigOverrides() want valid debug")
	}
	if !validateBoolFlag("dryrun", true, cfg.dryrun, t) {
		t.Error("TestConfigOverrides() want valid dryrun")
	}
}
