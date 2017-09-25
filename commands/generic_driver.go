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

// setup temp file for YAML --value parameter
func tempValues(r GenericDriver) (string, error) {
	file, err := ioutil.TempFile(os.TempDir(), "chartvalues")
	if err != nil {
		glog.Infof("err: %v\n", err)
		return "", err
	}
	glog.Infof("temporary file is %s", file.Name())

	if _, err := file.Write(r.YAMLValues); err != nil {
		glog.Infof("err: %v\n", err)
		os.Remove(file.Name())
		return "", err
	}
	if err := file.Close(); err != nil {
		glog.Infof("err: %v\n", err)
		os.Remove(file.Name())
		return "", err
	}
	return file.Name(), nil
}

// if credentials are present, they try login and return result
func registryLogin(r GenericDriver) ([]byte, error) {
	if r.Username != "" && r.Password != "" && r.Server != "" {
		// Login required for private application repos
		arguments := []string{"registry",
			"login",
			"-u " + r.Username,
			"-p " + r.Password,
			r.Server,
		}
		return r.execute(arguments)
	}
	return nil, nil
}

// convert internal json to yaml for api objects
func jsonToYaml(r *GenericDriver) error {
	glog.Infof("raw json: %s\n", r.JSONValues)
	json := []byte(r.JSONValues)
	yaml, err := yaml.JSONToYAML(json)
	if err != nil {
		glog.Infof("err: %v\n", err)
		return err
	}
	glog.Infof("raw yaml: %s\n", string(yaml))
	r.YAMLValues = yaml
	return nil
}

// Install - isntall the chart.
func (r GenericDriver) Install() ([]byte, error) {
	if output, err := registryLogin(r); err != nil {
		return output, err
	}

	filename, err := tempValues(r)
	if err != nil {
		return nil, err
	}
	defer os.Remove(filename)

	if err = jsonToYaml(&r); err != nil {
		return nil, err
	}

	// Do the install
	arguments := []string{"registry",
		"install",
		r.ChartLocation,
		"--namespace " + r.Namespace,
		"--name " + r.DeploymentName,
		"--values " + filename,
		"--version " + r.Version,
	}
	return r.execute(arguments)
}

// Upgrade - upgrade the chart.
func (r GenericDriver) Upgrade() ([]byte, error) {
	if output, err := registryLogin(r); err != nil {
		return output, err
	}

	filename, err := tempValues(r)
	if err != nil {
		return nil, err
	}
	defer os.Remove(filename)

	if err = jsonToYaml(&r); err != nil {
		return nil, err
	}

	// Do the upgrade
	arguments := []string{"registry",
		"upgrade",
		r.ChartLocation + "@" + r.Version,
		r.DeploymentName,
		"--values " + filename,
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
