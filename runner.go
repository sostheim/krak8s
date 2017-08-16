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
	"krak8s/commands"
	"krak8s/queue"
	"path"
	"sync"

	"github.com/golang/glog"
)

// Run back end tasks with queue semantics

// RequestType - requested tasks available
type RequestType int

const (
	// AddProject request
	AddProject RequestType = iota
	// RunningAddProject request
	RunningAddProject
	// UpdateProject request
	UpdateProject
	// RunningUpdateProject request
	RunningUpdateProject
	// RemoveProject request
	RemoveProject
	// RunningRemoveProject request
	RunningRemoveProject
	// AddChart request
	AddChart
	// UpdateChart request
	UpdateChart
	// RemoveChart request
	RemoveChart
)

func (req RequestType) String() string {
	return []string{
		"AddProject",
		"RunningAddProject",
		"UpdateProject",
		"RunningUpdateProject",
		"RemoveProject",
		"RunningRemoveProject",
		"AddChart",
		"UpdateChart",
		"RemoveChart",
	}[req]
}

// RequestStatus - request statuses
type RequestStatus int

const (
	// Waiting - request is queued, but not started processing yet
	Waiting RequestStatus = iota
	// Processing - request is being processed
	Processing
	// Deleting - a delete of the request is pending
	Deleting
	// Finished - request processing is complete
	Finished
	// Absent - request could not be found, it may have already finished or be deleted
	Absent
)

func (req RequestStatus) String() string {
	return []string{
		"Waiting",
		"Processing",
		"Deleting",
		"Finished",
		"Absent",
	}[req]
}

// Request - task to run
type Request struct {
	task        *queue.Task
	requestType RequestType
	name        string // either a project or chart name
	namespace   string
	nodes       int
	version     string
	registry    string
	config      string
	retry       int
}

var configFile string

// RunnerSetup - needs only be called once before execting ny funcitons below (note: this can't be func init()\)
func RunnerSetup() {
	if *krak8sCfg.krakenCommand == commands.K2 {
		commands.K2SetupEnv()
		configFile = path.Join(*krak8sCfg.krakenConfigDir, *krak8sCfg.krakenConfigFile)
	}
}

// NewRequest creates an request for processing
func NewRequest(req RequestType) *Request {
	return &Request{
		task:        queue.NewTask(),
		requestType: req,
	}
}

// Runner for request from API server to backend
type Runner struct {
	index           int
	pendingRequests map[int]*Request
	mutex           *sync.Mutex
	sync            chan int
}

// NewRunner creates a request runner
func NewRunner() *Runner {
	return &Runner{
		index:           0,
		pendingRequests: make(map[int]*Request),
		mutex:           &sync.Mutex{},
		sync:            make(chan int),
	}
}

// ProcessRequests - runner's main loop for request processing
func (r *Runner) ProcessRequests() {
	for {
		request := <-r.sync
		if request < 0 {
			return
		}
		r.handle(request)
	}
}

func (r *Runner) handle(index int) {
	done := false
	request := r.pendingRequests[index]
	if request.requestType >= AddProject && request.requestType <= RemoveProject {
		done = r.handleProjects(request)
	} else if request.requestType >= AddChart && request.requestType <= RemoveChart {
		done = r.handleCharts(request)
	}
	if done {
		r.DeleteRequest(index)
	}
}

func (r *Runner) handleProjects(request *Request) bool {

	// Check to see if the request is already (still) running
	if request.requestType == RunningAddProject ||
		request.requestType == RunningUpdateProject ||
		request.requestType == RunningRemoveProject {
		// Decrement the request type back to it's original value
		request.requestType--
		return true
	}

	cfg := commands.NewProjectConfig(request.name, request.nodes, request.namespace)
	cfg.KeyPair = *krak8sCfg.krakenKeyPair
	cfg.KubeConfigName = *krak8sCfg.krakenKubeConfig
	var command []string
	configPath := path.Join(*krak8sCfg.krakenConfigDir, *krak8sCfg.krakenConfigFile)
	if request.requestType == AddProject {
		err := commands.AddProjectTemplate(cfg, configPath)
		if err != nil {
			glog.Errorf("Discarding add: configuration update failure: %v", err)
			return true
		}
		if *krak8sCfg.krakenCommand == commands.K2 {
			command = commands.K2CmdUpdate(true, commands.K2ExtraVarsAddNodePools, *krak8sCfg.krakenConfigDir, configFile, request.name)
		} else {
			command = commands.ClusterUpdateAdd(request.name)
		}
	} else if request.requestType == RemoveProject {
		err := commands.DeleteProject(cfg, configPath)
		if err != nil {
			glog.Errorf("Discarding remove: configuration update failure: %v", err)
			return true
		}
		if *krak8sCfg.krakenCommand == commands.K2 {
			command = commands.K2CmdUpdate(true, commands.K2ExtraVarsRemoveNodePools, *krak8sCfg.krakenConfigDir, configFile, request.name)
		} else {
			command = commands.ClusterUpdateAdd(request.name)
		}
	} else {
		return true
	}

	// Block the command state in the queue and run the command to completion.
	queue.Started()
	var err error
	for err == nil && request.retry >= 0 {
		_, err = commands.Execute(command[0], command[1:])
		if err != nil {
			glog.Errorf("command execution retry count: %v", request.retry)
			glog.Errorf("command execution failed on: %v", err)
		}
		request.retry--
	}
	queue.Done()

	return true
}

func (r *Runner) handleCharts(request *Request) bool {
	return true
}

// DeleteRequest - remove request from processing pipeline
func (r *Runner) DeleteRequest(index int) {
	request, ok := r.pendingRequests[index]
	if !ok {
		return
	}

	// If the queued task is already (or still) running, it can't be deleted yet.
	status := queue.Delete(request.task.ID)
	if status == queue.Running {
		// Increment the request type indicate that it's running
		request.requestType--
		r.sync <- index
		return
	}

	glog.Infof("Queued task delted: type: %s, name: %s, namespace: %s, queueing duration: %s, running duration %s",
		request.requestType.String(), request.name, request.namespace, queue.QueuedDuration().String(), queue.RunningDuration().String())

	// ok to remove the request from the pending map
	r.mutex.Lock()
	defer r.mutex.Unlock()
	delete(r.pendingRequests, index)

	return
}

// ProjectRequest - submit project add request for processing.
func (r *Runner) ProjectRequest(action RequestType, name, namespace string, nodes int) RequestStatus {
	req := NewRequest(action)
	req.name = name
	req.namespace = namespace
	req.nodes = nodes
	req.retry = 1
	queue.Submit(req.task)

	// add the request to the pending map
	r.mutex.Lock()
	defer r.mutex.Unlock()
	r.index++
	r.pendingRequests[r.index] = req
	r.sync <- r.index

	return Waiting
}

// ProjectStatus - search for a request and, if found, return status
func (r *Runner) ProjectStatus(name, namespace string) RequestStatus {
	for _, req := range r.pendingRequests {
		if req.name == name && req.namespace == namespace {
			switch status := queue.Status(req.task.ID); status {
			case queue.Queued:
				return Waiting
			case queue.Running:
				return Processing
			}
		}
	}
	return Absent
}
