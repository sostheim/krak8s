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
	"io/ioutil"
	"os"

	"github.com/ghodss/yaml"
	"github.com/golang/glog"
)

// GenericDriver - control structure for deploying a generic chart.
type GenericDriver struct {
	DeploymentName string
	ChartLocation  string
	Version        string
	Server         string
	Namespace      string
	SetConfig      string
	JSONValues     string
	YAMLValues     []byte
	Username       string
	Password       string
}

func (r GenericDriver) jsonToYaml() {
	json := []byte(r.JSONValues)
	yaml, err := yaml.JSONToYAML(json)
	if err != nil {
		glog.Infof("err: %v\n", err)
		return
	}
	r.YAMLValues = yaml
}

// Install - isntall the chart.
func (r GenericDriver) Install() ([]byte, error) {
	file, err := ioutil.TempFile(os.TempDir(), "chartvalues")
	if err != nil {
		glog.Infof("err: %v\n", err)
		return nil, err
	}
	glog.Infof("temporary file is %s", file.Name())
	defer os.Remove(file.Name())

	if _, err := file.Write(r.YAMLValues); err != nil {
		glog.Infof("err: %v\n", err)
		return nil, err
	}
	if err := file.Close(); err != nil {
		glog.Infof("err: %v\n", err)
		return nil, err
	}

	// Login
	arguments := []string{"registry",
		"login",
		"-u " + r.Username,
		"-p " + r.Password,
		"quay.io",
	}
	r.execute(arguments)

	// Do the install
	arguments = []string{"registry",
		"install",
		r.ChartLocation,
		"--namespace " + r.Namespace,
		"--name " + r.DeploymentName,
		"--values " + file.Name(),
		"--version " + r.Version,
	}
	return r.execute(arguments)

}

// Upgrade - upgrade the chart.
func (r GenericDriver) Upgrade() ([]byte, error) {
	file, err := ioutil.TempFile(os.TempDir(), "chartvalues")
	if err != nil {
		glog.Infof("err: %v\n", err)
		return nil, err
	}
	glog.Infof("Template file is %s", file.Name())
	defer os.Remove(file.Name())

	if _, err := file.Write(r.YAMLValues); err != nil {
		glog.Infof("err: %v\n", err)
		return nil, err
	}

	if err := file.Close(); err != nil {
		glog.Infof("err: %v\n", err)
		return nil, err
	}

	// Login
	arguments := []string{"registry",
		"login",
		"-u " + r.Username,
		"-p " + r.Password,
		r.Server,
	}
	r.execute(arguments)

	// Do the upgrade
	arguments = []string{"registry",
		"upgrade",
		r.ChartLocation + "@0.1.0",
		r.DeploymentName,
		"--values " + file.Name(),
	}
	return r.execute(arguments)
}

// Remove - remove the chart.
func (r GenericDriver) Remove() ([]byte, error) {
	arguments := []string{"delete",
		"--purge",
		r.DeploymentName,
	}

	return r.execute(arguments)
}

func (r GenericDriver) execute(arguments []string) ([]byte, error) {
	return Execute("helm", arguments)
}
