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
	"text/template"
	"os"
	"log"
	"io/ioutil"
	"bytes"
	"os/exec"
)

const(
	MongoReplicasetTemplate = `scheduling:
  affinity:
    node:
      labels:
        - key: customer
          operator: In
          values: [ "{{ .CustomerName }}" ]
  tolerations:
    - key: customer
      value: {{ .CustomerName }}
      effect: NoSchedule
networkPolicy:
  ingress:
    enabled: true
    namespaceLabels:
      - key: customer
        value: {{ .CustomerName }}
    podLabels:
      - key: customer
        value: {{ .CustomerName }}
resources:
  limits:
    cpu: 200m
    memory: 512Mi
  requests:
    cpu: 200m
    memory: 512Mi`
)

type MongoReplicasetDriver struct {

	DeploymentName string
	ChartLocation string
	Namespace string

	CustomerName string

	Template string

}

func (m MongoReplicasetDriver) install() {
	templ, err := template.New("mongoTemplate").Parse(m.Template)
	if err != nil {
		log.Fatalf("execution failed: %s", err)
	}

	file, err := ioutil.TempFile(os.TempDir(), "prefix")
	log.Printf("Template file is %s", file.Name())
	defer os.Remove(file.Name())

	err = templ.Execute(file, m)

	arguments := []string{"registry",
		"install",
		m.ChartLocation,
		"--namespace " + m.Namespace,
		"--name " + m.DeploymentName,
		"--values " + file.Name(),
		"--version 1.2.0-0",
	}

	m.execute("/usr/local/bin/helm", arguments)

}


func (m MongoReplicasetDriver) upgrade() {
	templ, err := template.New("mongoTemplate").Parse(m.Template)
	if err != nil {
		log.Fatalf("execution failed: %s", err)
	}

	file, err := ioutil.TempFile(os.TempDir(), "prefix")
	log.Printf("Template file is %s", file.Name())
	defer os.Remove(file.Name())

	err = templ.Execute(file, m)


	arguments := []string{"registry",
		"upgrade",
		m.ChartLocation + "@1.2.0-0",
		m.DeploymentName,
		"--values " + file.Name(),
	}

	m.execute("/usr/local/bin/helm", arguments)
}

func (m MongoReplicasetDriver) remove() {
	arguments := []string{"delete",
		"--purge",
		m.DeploymentName,
	}

	m.execute("/usr/local/bin/helm", arguments)
}

func (m MongoReplicasetDriver) execute(command string, arguments []string) ([]byte, error) {
	cmd := exec.Command("helm", arguments...)
	stdoutBuf := &bytes.Buffer{}
	stderrBuf := &bytes.Buffer{}
	cmd.Stdout = stdoutBuf
	cmd.Stderr = stderrBuf

	if err := cmd.Run(); err != nil {
		log.Printf("k2cli.Execute(): cmd:  %s, args: %s returned error: %v", command, arguments, err)
		log.Printf("k2cli.Execute(): cmd:  %s, stderr: %s", command, string(stderrBuf.Bytes()))
		log.Printf("k2cli.Execute(): cmd:  %s, stdout: %v", command, string(stdoutBuf.Bytes()))
		return stderrBuf.Bytes(), err
	}
	return stdoutBuf.Bytes(), nil

}

/*
func main () {

	mongo := MongoReplicasetDriver{ DeploymentName: "test",
		ChartLocation: "quay.io/samsung_cnct/mongodb-replicaset",
		Namespace: "test",
		CustomerName: "joe",
		Template: MongoReplicasetTemplate,
	}

	//mongo.install()
	//mongo.upgrade
	//mongo.remove()

}
*/
