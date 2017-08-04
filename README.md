[![Go Report Card](https://goreportcard.com/badge/github.com/samsung-cnct/krak8s)](https://goreportcard.com/report/github.com/samsung-cnct/krak8s)
[![Docker Repository on Quay.io](https://quay.io/repository/samsung_cnct/krak8s/status "Docker Repository on Quay.io")](https://quay.io/repository/samsung_cnct/krak8s)
[![maturity](https://img.shields.io/badge/status-alpha-red.svg)](https://github.com/github.com/samsung-cnct/krak8s)

# API Service for Kraken and Kubernetes Commands
## Overview

A REST based API for managing [Kraken](https://github.com/samsung-cnct/k2) and [Kubernetes](https://kubernetes.io/) actions. is a Kubernetes Service Load balancer.

### Connectivity
A deployment of krak8s requires, at a minimum, network connectivity to both the Kubernetes API server and an execution environment for Kraken. The Kubernetes API server can be accessed via `kubectl proxy` for development, but this is not recommended for production deployments. For normal operation, the standard access via [`kubeconfig`](https://kubernetes.io/docs/concepts/cluster-administration/authenticate-across-clusters-kubeconfig/) or the Kubernetes API Server endpoint is supported.

## Running krak8s

As noted [above](#overview), krak8s can run as a static binary, a container image under a container runtime, or in a Kubernetes [Pod](https://kubernetes.io/docs/concepts/workloads/pods/pod/). Regardless of the runtime environment, krak8s has a number of command line options that define how it operates.
```
$ ./krak8s --help
Usage of ./krak8s:
      --alsologtostderr                  log to standard error as well as files
      --health-check                     enable health checking for API servcie (default false)
      --kubeconfig string                absolute path to the kubeconfig file
      --log_backtrace_at traceLocation   when logging hits line file:N, emit a stack trace (default :0)
      --log_dir string                   If non-empty, write log files in this directory
      --logtostderr                      log to standard error instead of files
      --proxy string                     kubctl proxy server running at the given url
      --service-name string              API Service name
      --stderrthreshold severity         logs at or above this threshold go to stderr (default 2)
  -v, --v Level                          log level for V logs
      --version                          display version info and exit
      --vmodule moduleSpec               comma-separated list of pattern=N settings for file-filtered logging
```
### Configuration Flags
Without going in to an explanation of all of the parameters, many of which should have sufficient explanation in the help provided, of particular interest to controlling the operation of krak8s are the following:<br />
<b>--health-check</b> - Defaults to true, but may be disabled by passing a value of false. Allows external service monitors to check the health of `krak8s` itself.<br />
<b>--kubeconfig</b> - Use the referenced kubeconfig for credentialed access to the cluster.<br />
<b>--proxy</b> - Use the `kubectl proxy` URL for access to the cluster. See for example [using kubectl proxy](https://kubernetes.io/docs/concepts/cluster-administration/access-cluster/#using-kubectl-proxy).<br />

### Environment Variables
krak8s is configurable through command line configuration flags, and through a subset of environment variables. Any configuration value set on the command line takes precedence over the same value from the environment.

The format of the environment variable for flag for flag is composed of the prefix `KRAK8S_` and the reamining text of the flag in all uppper case with all hyphens replaced by underscores.  Fore example, `--example-flag` would map to `KRAK8S_EXAMPLE_FLAG`. 

Not every flag can be set via an environment variable.  This is due to the fact that the set of flags is an aggregate of those that belong to krak8s and 3rd party Go packages.  The set of flags that do have corresponding environment variable support are listed below:
* --health-check
* --kubeconfig
* --service-name

### Details
The health check service is the HTTP endpoint `/healthz` by default...  

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
        args: ["--v=2", "--logtostderr=true"]
```

## Cutting a release

Install github-release from https://github.com/c4milo/github-release  
Create a github personal access token with repo read/write permissions and export it as GITHUB_TOKEN  
Adjust VERSION and TYPE variables in the [Makefile](Makefile) as needed  
Run ```make release```
