package kraken

import (
	"text/template"
	"os"
	"log"
	"io/ioutil"
	"bytes"
	"os/exec"
)

const(
	ReactionTemplate = `image: {{.Image}}
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

type ReactionDriver struct {

	DeploymentName string
	ChartLocation string
	Namespace string

	Image string
	ImageTag string
	ImagePullSecret string
	DefaultHostName string
	MainHostName string
	RootHostName string
	CustomerName string

	Username string
	Password string

	Template string

}

func (r ReactionDriver) install() {
	templ, err := template.New("mongoTemplate").Parse(r.Template)
	if err != nil {
		log.Fatalf("execution failed: %s", err)
	}

	file, err := ioutil.TempFile(os.TempDir(), "prefix")
	log.Printf("Template file is %s", file.Name())
	defer os.Remove(file.Name())

	err = templ.Execute(file, r)

	// Login
	arguments := []string{"registry",
		"login",
		"-u " + r.Username,
		"-p " + r.Password,
		"quay.io",
	}
	r.execute("/usr/local/bin/helm", arguments)

	// Do the install
	arguments = []string{"registry",
		"install",
		r.ChartLocation,
		"--namespace " + r.Namespace,
		"--name " + r.DeploymentName,
		"--values " + file.Name(),
		"--version 0.1.0",
	}
	r.execute("/usr/local/bin/helm", arguments)

}


func (r ReactionDriver) upgrade() {
	templ, err := template.New("mongoTemplate").Parse(r.Template)
	if err != nil {
		log.Fatalf("execution failed: %s", err)
	}

	file, err := ioutil.TempFile(os.TempDir(), "prefix")
	log.Printf("Template file is %s", file.Name())
	defer os.Remove(file.Name())

	err = templ.Execute(file, r)


	// Login
	arguments := []string{"registry",
		"login",
		"-u " + r.Username,
		"-p " + r.Password,
		"quay.io",
	}
	r.execute("/usr/local/bin/helm", arguments)

	// Do the upgrade
	arguments = []string{"registry",
		"upgrade",
		r.ChartLocation + "@0.1.0",
		r.DeploymentName,
		"--values " + file.Name(),
	}
	r.execute("/usr/local/bin/helm", arguments)
}

func (r ReactionDriver) remove() {
	arguments := []string{"delete",
		"--purge",
		r.DeploymentName,
	}

	r.execute("/usr/local/bin/helm", arguments)
}

func (r ReactionDriver) execute(command string, arguments []string) ([]byte, error) {
	cmd := exec.Command(command, arguments...)
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

	reaction := ReactionDriver{ DeploymentName: "test",
		ChartLocation: "quay.io/reactioncommerce/reactioncommerce",
		Namespace: "test",
		Image: "reactioncommerce/reaction",
		ImageTag: "latest",
		ImagePullSecret: "",
		DefaultHostName: "test.getreaction.io",
		MainHostName: "www.test.io",
		RootHostName: "www.test.io",
		CustomerName: "joe",
		Username: "reactioncommerce+reactioncommercero",
		Password: "somethingsomething",
		Template: ReactionTemplate,
	}

	reaction.install()
	//reaction.upgrade()
	//reaction.remove()
}

*/