[![Go Report Card](https://goreportcard.com/badge/github.com/samsung-cnct/krak8s)](https://goreportcard.com/report/github.com/samsung-cnct/krak8s)
[![Docker Repository on Quay.io](https://quay.io/repository/samsung_cnct/krak8s/status "Docker Repository on Quay.io")](https://quay.io/repository/samsung_cnct/krak8s)
[![maturity](https://img.shields.io/badge/status-alpha-red.svg)](https://github.com/github.com/samsung-cnct/krak8s)

# API Service for Kraken and Kubernetes Commands
## Overview

A REST based API for managing [Kraken](https://github.com/samsung-cnct/k2) and [Kubernetes](https://kubernetes.io/) actions.

Go directly to the [API Documentation and Specification](https://github.com/samsung-cnct/krak8s/blob/master/API%20Definitions.md) section.

### Connectivity
A deployment of krak8s requires network connectivity to the Kubernetes API server. The Kubernetes API server can be accessed via `kubectl proxy` for development, but this is not recommended for production deployments. For normal operation, the standard access via [`kubeconfig`](https://kubernetes.io/docs/concepts/cluster-administration/authenticate-across-clusters-kubeconfig/) or the Kubernetes API Server endpoint is supported.

## Running krak8s

The krak8s service can run as a static binary, a container image under a container runtime, or in a Kubernetes [Pod](https://kubernetes.io/docs/concepts/workloads/pods/pod/). Regardless of the runtime environment, krak8s has a number of command line options that define how it operates.
```
$ ./krak8s --help
Usage of ./krak8s:
      --alsologtostderr                  log to standard error as well as files
      --debug                            enable debug output
      --dry-run                          don't actually execute backend commands
      --health-check                     enable health checking for API service
      --kraken-command k2                command to run to execute kraken operations, either k2, or `k2cli` only (default "k2")
      --kraken-config-dir string         kraken configuration yaml directory path (default "${HOME}/.kraken")
      --kraken-config-file string        kraken configuration yaml file name (default "config.yaml")
      --kraken-kubeconfig string         kraken confiuration yaml: deployment.clusters[0].nodePools.kubeConfig (default "defaultKube")
      --kraken-nodepool-keypair string   kraken configuration yaml: deployment.clusters[0].nodePools.keyPair (default "defaultKeyPair")
      --kubeconfig string                absolute path to the kubeconfig file
      --log_backtrace_at traceLocation   when logging hits line file:N, emit a stack trace (default :0)
      --log_dir string                   If non-empty, write log files in this directory
      --logtostderr                      log to standard error instead of files
      --proxy string                     kubctl proxy server running at the given url
      --stderrthreshold severity         logs at or above this threshold go to stderr (default 2)
  -v, --v Level                          log level for V logs
      --version                          display version info and exit
      --vmodule moduleSpec               comma-separated list of pattern=N settings for file-filtered logging
```
### Configuration Flags
Without going into an explanation of all of the parameters, many of which should have sufficient explanation in the help provided, of particular interest to controlling the operation of krak8s are the following:<br />
<b>--debug</b> - Allow generation of additional output for debugging purposes.<br />
<b>--dry-run</b> - Prevent any backend services from being executed against the live cluster.<br />
<b>--health-check</b> - Allows external service monitors to check the health of the `krak8s` service.<br />
<b>--kubeconfig</b> - Use the referenced kubeconfig for credentialed access to the cluster.<br />
<b>--proxy</b> - Use the `kubectl proxy` URL for access to the cluster. See for example [using kubectl proxy](https://kubernetes.io/docs/concepts/cluster-administration/access-cluster/#using-kubectl-proxy).<br />
<b>--kraken-command</b> - The command to run to execute kraken operations, this can only be either `k2`, or `k2cli`<br />
<b>--kraken-config-dir</b> - The Kraken configuration yaml directory path (default "${HOME}/.kraken")<br />
<b>--kraken-config-file</b> - The Kraken configuration yaml file name (default "config.yaml")<br />
<b>--kraken-kubeconfig</b> - Value for Kraken confiuration yaml: deployment.clusters[0].nodePools.kubeConfig (default "defaultKube")<br />
<b>--kraken-nodepool-keypair</b> - Value for Kraken configuration yaml: deployment.clusters[0].nodePools.keyPair (default "defaultKeyPair")<br />

### Environment Variables
krak8s is configurable through command line configuration flags, and through a subset of environment variables. Any configuration value set on the command line takes precedence over the same value from the environment.

The format of the environment variable for a flag is composed of the prefix `KRAK8S_` and the remaining text of the flag in all uppercase with all hyphens replaced by underscores.  Fore example, `--example-flag` would map to `KRAK8S_EXAMPLE_FLAG`. 

Not every flag can be set via an environment variable.  This is due to the fact that the set of flags is an aggregate of those that belong to krak8s and 3rd party Go packages.  The set of flags that do have corresponding environment variable support are listed below:
* --debug
* --dry-run
* --health-check
* --kraken-command
* --kraken-config-dir
* --kraken-config-file
* --kraken-kubeconfig
* --kraken-nodepool-keypair
* --kubeconfig
* --proxy

### Details
The health check service HTTP endpoint is available at: `/healthz`.  

## Deploying krak8s Example
The following is an example of using a Kubernetes Deployment to run krak8s. 
```
aapiVersion: extensions/v1beta1
kind: Deployment
metadata:
  name: krak8s
  labels:
    name: krak8s
    app: krak8s
    version: 0.1.0
spec:
  replicas: 1
  template:
    metadata:
      labels:
        app: krak8s
        version: 0.1.0
    spec:
      containers:
      - name: krak8s
        image: quay.io/samsung_cnct/krak8s:latest
        args: ["--v=2", "--logtostderr=true", "--kraken-kubeconfig meteorKube", "--kraken-nodepool-keypair meteorKeyPair"]
```

## Building and Running the Project Locally
The following steps assume that you have a working Golang development environment on your local machine or in a container that you use for the same purpose.

1. Clone the Repository

The [Goa Design](https://github.com/goadesign/goa) framework for building micro-services and REST APIs has a very specific path requirement for projects based on it.  In particular, a Goa project must be based at `$GOPATH\src\<project-name>` for the code generator to work properly.

```
$ echo $GOPATH
/Users/sostheim/work

$ git clone https://github.com/samsung-cnct/krak8s.git ${GOPATH}/src/krak8s
Cloning into '/Users/sostheim/work/src/krak8s'...
remote: Counting objects: 2870, done.
remote: Compressing objects: 100% (320/320), done.
remote: Total 2870 (delta 183), reused 324 (delta 99), pack-reused 2438
Receiving objects: 100% (2870/2870), 6.91 MiB | 1.23 MiB/s, done.
Resolving deltas: 100% (1282/1282), done.
```

2. Project Make Targets

The tools necessary to build and deploy the project are all bundled in the Makefile.

  1. Make Dependencies

  The project has only two dependencies to build and run (cutting a release is discussed separately below).  These are the [Gox](https://github.com/mitchellh/gox) Go Cross Compiler, and the [Godep](https://github.com/tools/godep) dependency tool for go.
  ```
  $ make deps
  go get github.com/mitchellh/gox
  go get github.com/tools/godep
  ```
  The project maintains a `/vendor` directory of all the other dependencies of the project pinned to the requisite version.

  2. Make the Project

  The default target for the Makefile is the `build:` target.  This is a convenience to prevent building an unwanted target.  You should, as a general practices, always specify the build target on the command line.

  ```
  $ make build
  go build -ldflags "-X main.MajorMinorPatch=0.2.91 \
      -X main.ReleaseType=alpha \
      -X main.GitCommit=fbcb8984a595eb8e563992392d2e088339ca6aab"
  ```
  **Note** that the build automatically pulls the last Git Commit SHA-1 to tag the build output with, along with a [SemVer](http://semver.org/) version number and release type that are populated as defined in the Makefile itself.

That's it.  The project can now be run locally with a command string similar to the following.
```
./krak8s --alsologtostderr --kubeconfig /Users/sostheim/.kube/config --kraken-kubeconfig meteorKube --kraken-nodepool-keypair meteorKeyPair --kraken-config-file config.yaml --kraken-config-dir /Users/sostheim/.kraken --debug --dry-run
```
In the example above, several of the default values are specified on the command line, and could be omitted.  They are shown here for completeness.

**NOTE:** In this example, we have supplied both the `--debug` and `--dry-run` flags.  This is solely for the purpose of development and integration testing.  These flags **must** be removed to run the API service in a production environment where real changes to an active Kubernetes cluster are going to be performed.  As noted above in the description of the [Configuration Flags](#configuration-flags), the presence of the debug and dry run flags allow additional output to be generated, and allows the API to be exercised without affecting the live cluster.

## Building and Running the Project Container - MacOS
As before, the following steps also assume that you have a working Golang development environment on your local machine or in a container that you use for the same purpose.

1. Project Make Target: Push

As before, the tools necessary to build and deploy the project are all bundled in the Makefile.

  1. Make Push

  The command below will cross compile for the `amd64` target and push the image to the project's [CNCT Quay Repository](https://quay.io/repository/samsung_cnct/krak8s).
  ```
  $ make push
  go get github.com/mitchellh/gox
  go get github.com/tools/godep
  gox -ldflags "-X main.MajorMinorPatch=0.2.91 \
      -X main.ReleaseType=alpha \
      -X main.GitCommit=fbcb8984a595eb8e563992392d2e088339ca6aab -w" \
    -osarch="linux/amd64" \
    -output "build/{{.OS}}_{{.Arch}}/krak8s" \
    ./...
  Number of parallel builds: 7

  -->     linux/amd64: krak8s
  -->     linux/amd64: krak8s/tool/krak8s-cli
  docker build --rm --pull --tag quay.io/samsung_cnct/krak8s:latest .
  Sending build context to Docker daemon  108.3MB
  Step 1/8 : FROM quay.io/samsung_cnct/k2:latest
  latest: Pulling from samsung_cnct/k2

  [ . . . Clipped Command Output Text ... ]

  latest: digest: sha256:98bbaa5d2139616aa92913d883607584111e2330898ab69a160e7092537cb65d size: 44545
  ```
  **Note 1:** When building for yourself, you will need to substitute a Docker Image Repository that you have write access to in place of the Samsung Quay.io repository.

  **Note 2:** This is *NOT* a fully statically linked binary.  There are shared library dependencies that must be resolved on the target system when the binary runs.  This is intentional to allow the container to run locally on your Docker for Mac for testing.

## Building and Running the Project Container - Linux
As before, the following steps also assume that you have a working Golang development environment on your local machine or in a container that you use for the same purpose.

1. Project Make Target: Push Static

As before, the tools necessary to build and deploy the project are all bundled in the Makefile.  Here however, the build tool chain assumes that the build, tag, and push are all being executed on the same architecture as the deployment target architecture.  In this case that architecture is `amd64`.  As such we can pass the flags to statically link against the system libraries for deployment as standalone binary.  All other elements remain the same.

  1. Make Push Static

  As before, the command shown below will cross compile for the `amd64` target and push the image to the project's [CNCT Quay Repository](https://quay.io/repository/samsung_cnct/krak8s).
  ```
  $ make push_static
  go get github.com/mitchellh/gox
  go get github.com/tools/godep
  gox -ldflags "-X main.MajorMinorPatch=0.2.91 \
      -X main.ReleaseType=alpha \
      -X main.GitCommit=fbcb8984a595eb8e563992392d2e088339ca6aab -w \
      -linkmode external -extldflags -static" \
    -osarch="linux/amd64" \
    -output "build/{{.OS}}_{{.Arch}}/krak8s" \
    ./...
  Number of parallel builds: 7

  -->     linux/amd64: krak8s
  -->     linux/amd64: krak8s/tool/krak8s-cli
  docker build --rm --pull --tag quay.io/samsung_cnct/krak8s:latest .
  Sending build context to Docker daemon  108.3MB
  Step 1/8 : FROM quay.io/samsung_cnct/k2:latest
  latest: Pulling from samsung_cnct/k2

  [ . . . Clipped Command Output Text ... ]

  latest: digest: sha256:98bbaa5d2139616aa92913d883607584111e2330898ab69a160e7092537cb65d size: 44545
  ```
  **Note 3:** As before, when building for yourself, you will need to substitute a Docker Image Repository that you have write access to in place of the Samsung Quay.io repository.

## Cutting a release

Install github-release from https://github.com/c4milo/github-release  
Create a github personal access token with repo read/write permissions and export it as GITHUB_TOKEN  
Adjust VERSION and TYPE variables in the [Makefile](Makefile) as needed  
Run ```make release```
