package kraken

import (
	"bytes"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"text/template"
	"time"

	"github.com/golang/glog"
)

const (
	// DefautlKubeConfig - defautl kubeConfig in ~/.kraken/config.yaml::kubeConfigs
	DefautlKubeConfig = "defaultKube"
	// DefautlKeyPair - defautl keyPair in ~/.kraken/config.yaml::keyPairs
	DefautlKeyPair = "defaultKeyPair"
	// NodePoolMarker - generated node pool configuration insertion point
	NodePoolMarker = "# |--> NODE_POOL_MARKER <--|"
	// ServicesMarker - generated services configuration insertion point
	ServicesMarker = "# |--> SERVICES_MARKER <--|"
)

// ProjectConfig describes the cluster resource configuration for a project
type ProjectConfig struct {
	Name           string
	NodePoolCount  int
	KubeConfigName string
	KeyPair        string
	Namespace      string
}

// NewProjectConfig creates an configuration record with default values.
func NewProjectConfig(name string, count int, ns string) ProjectConfig {
	return ProjectConfig{
		Name:           name,
		NodePoolCount:  count,
		Namespace:      ns,
		KeyPair:        DefautlKeyPair,
		KubeConfigName: DefautlKubeConfig,
	}
}

func copyFile(dst, src string, perm os.FileMode) error {
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()
	tmp, err := ioutil.TempFile(filepath.Dir(dst), "")
	if err != nil {
		return err
	}
	_, err = io.Copy(tmp, in)
	if err != nil {
		tmp.Close()
		os.Remove(tmp.Name())
		return err
	}
	if err = tmp.Close(); err != nil {
		os.Remove(tmp.Name())
		return err
	}
	if err = os.Chmod(tmp.Name(), perm); err != nil {
		os.Remove(tmp.Name())
		return err
	}
	return os.Rename(tmp.Name(), dst)
}

func copyConfigFileBackup(path string) error {
	src := path
	dest := path + "." + strconv.FormatInt(time.Now().Unix(), 10)
	if _, err := os.Stat(src); os.IsNotExist(err) {
		return err
	}
	if err := copyFile(dest, src, os.FileMode(0644)); err != nil {
		return err
	}
	return nil
}

func templateNodes(config ProjectConfig) (string, error) {
	// path.Join(configPath, "config.yaml"")
	tmpl, err := template.New("node_pool.tmpl").ParseFiles("node_pool.tmpl")
	if err != nil {
		glog.Warningf("failed to parse node pool template file: %v", err)
		return "", err
	}
	buf := new(bytes.Buffer)
	if err := tmpl.Execute(buf, config); err != nil {
		glog.Warningf("failed to execute node pool template %v", err)
		return "", err
	}
	return buf.String(), nil
}

func templateServices(config ProjectConfig) (string, error) {
	// path.Join(configPath, "config.yaml"")
	tmpl, err := template.New("services.tmpl").ParseFiles("services.tmpl")
	if err != nil {
		glog.Warningf("failed to parse services template file: %v", err)
		return "", err
	}
	buf := new(bytes.Buffer)
	if err := tmpl.Execute(buf, config); err != nil {
		glog.Warningf("failed to execute services template %v", err)
		return "", err
	}
	return buf.String(), nil
}

// AddProjectTemplate - copies the configuration file, which *MUST* be the
// most current up to date configuration file, and then inserts a node pool
// and service stanza in to the configuration file for the project.
func AddProjectTemplate(config ProjectConfig, filename string) error {

	nodeConfig, err := templateNodes(config)
	if err != nil {
		glog.Warning("failed to template nodePool config stanza")
		return err
	}

	svcConfig, err := templateServices(config)
	if err != nil {
		glog.Warning("failed to template services config stanza")
		return err
	}

	err = copyConfigFileBackup(filename)
	if err != nil {
		glog.Warning("failed to make backup copy of config file")
		return err
	}

	configFileData, err := ioutil.ReadFile(filename)
	if err != nil {
		glog.Warning("unable to open config file")
		return err
	}

	configFileLines := strings.Split(string(configFileData), "\n")

	var outputFileLines []string
	for _, line := range configFileLines {
		if strings.Contains(line, config.Name) || strings.Contains(line, config.Namespace) {
			glog.Infof("Configuration file already contains the Project Name: %s, or Namespace: %s",
				config.Name, config.Name)
			return nil
		}
		outputFileLines = append(outputFileLines, line)
		if strings.Contains(line, NodePoolMarker) {
			outputFileLines = append(outputFileLines, nodeConfig)
		} else if strings.Contains(line, ServicesMarker) {
			outputFileLines = append(outputFileLines, svcConfig)
		}
	}

	outputFileData := strings.Join(outputFileLines, "\n")
	err = ioutil.WriteFile(filename, []byte(outputFileData), 0644)
	if err != nil {
		glog.Warning("failed writing out new version of config file")
		return err
	}

	return nil
}

// DeleteProject - copies the configuration file, which *MUST* be the most
// current up to date configuration file, and then searches for and removes
// a node pool and service stanza from the configuration file for the project.
func DeleteProject(config ProjectConfig, filename string) error {

	err := copyConfigFileBackup(filename)
	if err != nil {
		glog.Warning("failed to make backup copy of config file")
		return err
	}

	configFileData, err := ioutil.ReadFile(filename)
	if err != nil {
		glog.Warning("unable to open config file")
		return err
	}

	configFileLines := strings.Split(string(configFileData), "\n")

	skip := 0
	var outputFileLines []string
	for _, line := range configFileLines {
		if strings.Contains(line, config.Name) {
			if strings.Contains(line, config.Name+"Nodes") {
				skip = 12
			} else if strings.Contains(line, config.Name+"-mongodb") {
				skip = 5
			}
		}
		if skip > 0 {
			skip--
			continue
		}
		outputFileLines = append(outputFileLines, line)
	}

	outputFileData := strings.Join(outputFileLines, "\n")
	err = ioutil.WriteFile(filename, []byte(outputFileData), 0644)
	if err != nil {
		glog.Warning("failed writing out new version of config file")
		return err
	}

	return nil
}
