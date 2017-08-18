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
	"log"
	"os"
	"text/template"

	"github.com/golang/glog"
)

const (
	// GenericTemplate - chart template
	GenericTemplate = `image: {{.Image}}
imageTag: {{.ImageTag}}
imagePullSecret: {{.ImagePullSecret}}
ingress:
  defaultHost:
    hostname: {{.DefaultHostName}}
  mainHost:
    hostname: {{.MainHostName}}
rootURL: https://{{.RootHostName}}
mongo:
  application:
    deploymentName: {{.CustomerName}}-mongo
  opLog:
    deploymentName: {{.CustomerName}}-mongo
scheduling:
  affinity:
    node:
      labels:
        - key: customer
          operator: In
          values: [ "{{.CustomerName}}" ]
  tolerations:
    - key: customer
      value: {{.CustomerName}}
      effect: NoSchedule
resources:
  limits:
    cpu: 800m
    memory: 1536Mi
  requests:
    cpu: 800m
    memory: 1536Mi`
)

// GenericDriver - control structure for deploying a generic chart.
type GenericDriver struct {
	DeploymentName string
	ChartLocation  string
	Namespace      string

	Image           string
	ImageTag        string
	ImagePullSecret string
	DefaultHostName string
	MainHostName    string
	RootHostName    string
	CustomerName    string

	Username string
	Password string

	Template string
}

// Install - isntall the chart.
func (r GenericDriver) Install() ([]byte, error) {
	templ, err := template.New("mongoTemplate").Parse(r.Template)
	if err != nil {
		log.Fatalf("execution failed: %s", err)
	}

	file, err := ioutil.TempFile(os.TempDir(), "prefix")
	log.Printf("Template file is %s", file.Name())
	defer os.Remove(file.Name())

	err = templ.Execute(file, r)
	if err != nil {
		glog.Warningf("failed to parse chart template: %v", err)
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
		"--version 0.1.0",
	}
	return r.execute(arguments)

}

// Upgrade - upgrade the chart.
func (r GenericDriver) Upgrade() ([]byte, error) {
	templ, err := template.New("mongoTemplate").Parse(r.Template)
	if err != nil {
		log.Fatalf("execution failed: %s", err)
	}

	file, err := ioutil.TempFile(os.TempDir(), "prefix")
	log.Printf("Template file is %s", file.Name())
	defer os.Remove(file.Name())

	err = templ.Execute(file, r)
	if err != nil {
		glog.Warningf("failed to parse chart template: %v", err)
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
