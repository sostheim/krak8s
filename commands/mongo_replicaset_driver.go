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
	// MongoReplicasetTemplate - template for mongo replica set.
	MongoReplicasetTemplate = `scheduling:
  affinity:
    node:
      labels:
        - key: nodepool
          operator: In
          values: [ "{{ .CustomerName }}Nodes" ]
  tolerations:
    - key: customer
      value: {{ .CustomerName }}
      effect: NoSchedule
resources:
  limits:
    cpu: 200m
    memory: 512Mi
  requests:
    cpu: 200m
    memory: 512Mi`
)

// MongoReplicasetDriver - control structure for deploying mongo replica set.
type MongoReplicasetDriver struct {
	DeploymentName string
	ChartLocation  string
	Namespace      string

	CustomerName string

	Template string
}

// Install - upgrade the mongo replicaset chart.
func (m MongoReplicasetDriver) Install() ([]byte, error) {
	templ, err := template.New("mongoTemplate").Parse(m.Template)
	if err != nil {
		log.Fatalf("execution failed: %s", err)
	}

	file, err := ioutil.TempFile(os.TempDir(), "prefix")
	log.Printf("Template file is %s", file.Name())
	defer os.Remove(file.Name())

	err = templ.Execute(file, m)
	if err != nil {
		glog.Warningf("failed to parse mongo template: %v", err)
		return nil, err
	}

	arguments := []string{"registry",
		"install",
		m.ChartLocation,
		"--namespace " + m.Namespace,
		"--name " + m.DeploymentName,
		"--values " + file.Name(),
		"--version 1.2.0-0",
	}

	return m.execute(arguments)
}

// Upgrade - upgrade the mongo replicaset chart.
func (m MongoReplicasetDriver) Upgrade() ([]byte, error) {
	templ, err := template.New("mongoTemplate").Parse(m.Template)
	if err != nil {
		log.Fatalf("execution failed: %s", err)
	}

	file, err := ioutil.TempFile(os.TempDir(), "prefix")
	log.Printf("Template file is %s", file.Name())
	defer os.Remove(file.Name())

	err = templ.Execute(file, m)
	if err != nil {
		glog.Warningf("failed to parse mongo template: %v", err)
		return nil, err
	}

	arguments := []string{"registry",
		"upgrade",
		m.ChartLocation + "@1.2.0-0",
		m.DeploymentName,
		"--values " + file.Name(),
	}

	return m.execute(arguments)
}

// Remove - remove the mongo replicaset chart.
func (m MongoReplicasetDriver) Remove() ([]byte, error) {
	arguments := []string{"delete",
		"--purge",
		m.DeploymentName,
	}

	return m.execute(arguments)
}

func (m MongoReplicasetDriver) execute(arguments []string) ([]byte, error) {
	return Execute("helm", arguments)
}
