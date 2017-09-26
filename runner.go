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
	"os"
	"path"
	"sync"
	"time"

	"github.com/golang/glog"
)

// Run back end tasks with queue semantics

var configFile string

// RunnerSetup - needs only be called once before execting ny functions below (note: this can't be func init()\)
func RunnerSetup() {
	if *krak8sCfg.krakenCommand == commands.K2 {
		commands.K2SetupEnv()
		configFile = path.Join(*krak8sCfg.krakenConfigDir, *krak8sCfg.krakenConfigFile)
	}
}

// RequestType - requested tasks available
type RequestType int

const (
	// AddProject request
	AddProject RequestType = iota
	// UpdateProject request
	UpdateProject
	// RemoveProject request
	RemoveProject
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
	dataStore   *DataStore
	projObj     *ProjectObject
	nsObj       *NamespaceObject
	resObj      *ResourceObject
	appObj      *ApplicationObject
	retryCount  int
}

// NewResourceRequest creates an request for processing
func NewResourceRequest(req RequestType, ds *DataStore, proj *ProjectObject, ns *NamespaceObject, obj *ResourceObject) *Request {
	return &Request{
		task:        queue.NewTask(),
		dataStore:   ds,
		projObj:     proj,
		nsObj:       ns,
		resObj:      obj,
		requestType: req,
	}
}

// NewChartRequest creates an request for processing
func NewChartRequest(req RequestType, ds *DataStore, proj *ProjectObject, ns *NamespaceObject, app *ApplicationObject) *Request {
	return &Request{
		task:        queue.NewTask(),
		dataStore:   ds,
		projObj:     proj,
		nsObj:       ns,
		appObj:      app,
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

	cfg := commands.NewProjectConfig(request.projObj.Name, request.resObj.NodePoolSize, request.nsObj.Name)
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
			command = commands.K2CmdUpdate(*krak8sCfg.krakenInDocker, commands.K2ExtraVarsAddNodePools, *krak8sCfg.krakenConfigDir, configFile, request.projObj.Name)
		} else {
			command = commands.ClusterUpdateAdd(request.projObj.Name)
		}

		request.resObj.State = ResourceStarting
		request.resObj.UpdatedAt = time.Now()

	} else if request.requestType == RemoveProject {

		err := commands.DeleteProject(cfg, configPath)
		if err != nil {
			glog.Errorf("Discarding remove: configuration update failure: %v", err)
			return true
		}

		if *krak8sCfg.krakenCommand == commands.K2 {
			command = commands.K2CmdUpdate(*krak8sCfg.krakenInDocker, commands.K2ExtraVarsRemoveNodePools, *krak8sCfg.krakenConfigDir, configFile, request.projObj.Name)
		} else {
			command = commands.ClusterUpdateAdd(request.projObj.Name)
		}

		request.resObj.State = ResourceDeleting
		request.resObj.UpdatedAt = time.Now()

	} else {
		return true
	}

	runProjectRequestWithRetries(request, command)

	return true
}

func runProjectRequestWithRetries(request *Request, command []string) {
	// Block the command state in the queue and run the command to completion.
	queue.Started()
	if *krak8sCfg.krakenInDocker == false {
		if wd, err := os.Getwd(); err == nil {
			os.Chdir("/kraken")
			defer os.Chdir(wd)
		} else {
			glog.Infof("error getting current working directory: %v", err)
		}
	}

	tries := request.retryCount
	for tries >= 0 {
		_, err := commands.Execute(command[0], command[1:])
		if err != nil {
			tries--
			glog.Errorf("command execution retry count: %v", request.retryCount)
			glog.Errorf("command execution failed on: %v", err)
			if request.resObj.State == ResourceCreateRequested || request.resObj.State == ResourceStarting {
				request.resObj.State = ResourceErrorStarting
			} else if request.resObj.State == ResourceDeleteRequested || request.resObj.State == ResourceDeleting {
				request.resObj.State = ResourceErrorDeleting
			}
		} else {
			if *krak8sCfg.debug {
				glog.Infof("command execution success, tries: %d", tries)
			}
			tries = -1
			if request.resObj.State == ResourceCreateRequested || request.resObj.State == ResourceStarting || request.resObj.State == ResourceErrorStarting {
				request.resObj.State = ResourceActive
			} else if request.resObj.State == ResourceDeleteRequested || request.resObj.State == ResourceDeleting || request.resObj.State == ResourceErrorDeleting {
				request.resObj.State = ResourceDeleted
			}
		}
		request.resObj.UpdatedAt = time.Now()
	}
	queue.Done()
}

func (r *Runner) handleCharts(request *Request) bool {
	if request.appObj.ChartName == "mongodb-replicaset" {
		return r.handleMongoChart(request)
	}
	return r.handleGenericChart(request)
}

func (r *Runner) handleMongoChart(request *Request) bool {

	mongo := commands.MongoReplicasetDriver{
		DeploymentName: request.projObj.Name + "-mongodb",
		ChartLocation:  request.appObj.Server + "/" + request.appObj.ChartRegistry + "/" + request.appObj.ChartName,
		Namespace:      request.nsObj.Name,
		CustomerName:   request.projObj.Name,
		Template:       commands.MongoReplicasetTemplate,
	}
	if request.requestType == AddChart {
		request.appObj.Status.State = ApplicationUnknown
		request.appObj.UpdatedAt = time.Now()
		// Block the command state in the queue and run the command to completion.
		queue.Started()
		tries := request.retryCount
		for tries >= 0 {
			_, err := mongo.Install()
			if err != nil {
				tries--
				glog.Errorf("install retry count: %v", request.retryCount)
				glog.Errorf("install failed on: %v", err)
				request.appObj.Status.State = ApplicationFailed
			} else {
				if *krak8sCfg.debug {
					glog.Infof("command execution success, tries: %d", tries)
				}
				tries = -1
				request.appObj.Status.State = ApplicationDeployed
				request.appObj.Status.DeployedAt = time.Now()
			}
			request.appObj.UpdatedAt = time.Now()
		}
		queue.Done()
	} else if request.requestType == RemoveChart {
		request.appObj.Status.State = ApplicationDeleting
		request.appObj.UpdatedAt = time.Now()
		// Block the command state in the queue and run the command to completion.
		queue.Started()
		tries := request.retryCount
		for tries >= 0 {
			_, err := mongo.Remove()
			if err != nil {
				tries--
				glog.Errorf("remove retry count: %v", request.retryCount)
				glog.Errorf("remove failed on: %v", err)
				request.appObj.Status.State = ApplicationFailed
			} else {
				if *krak8sCfg.debug {
					glog.Infof("command execution success, tries: %d", tries)
				}
				tries = -1
				request.appObj.Status.State = ApplicationDeleted
				request.appObj.Status.DeployedAt = time.Now()
			}
			request.appObj.UpdatedAt = time.Now()
		}
		queue.Done()
	}

	return true
}

func (r *Runner) handleGenericChart(request *Request) bool {
	chart := commands.GenericDriver{
		DeploymentName: request.appObj.Deployment,
		ChartLocation:  request.appObj.Server + "/" + request.appObj.ChartRegistry + "/" + request.appObj.ChartName,
		Version:        request.appObj.ChartVersion,
		Server:         request.appObj.Server,
		SetConfig:      request.appObj.Config,
		JSONValues:     request.appObj.JSONValues,
		Namespace:      request.nsObj.Name,
		Username:       request.appObj.Username,
		Password:       request.appObj.Password,
	}

	if request.requestType == AddChart {
		request.appObj.Status.State = ApplicationUnknown
		request.appObj.UpdatedAt = time.Now()
		// Block the command state in the queue and run the command to completion.
		queue.Started()
		tries := request.retryCount
		for tries >= 0 {
			_, err := chart.Install()
			if err != nil {
				tries--
				glog.Errorf("install retry count: %v", request.retryCount)
				glog.Errorf("install failed on: %v", err)
				request.appObj.Status.State = ApplicationFailed
			} else {
				if *krak8sCfg.debug {
					glog.Infof("command execution success, tries: %d", tries)
				}
				tries = -1
				request.appObj.Status.State = ApplicationDeployed
				request.appObj.Status.DeployedAt = time.Now()
			}
			request.appObj.UpdatedAt = time.Now()
		}
		queue.Done()
	} else if request.requestType == RemoveChart {
		request.appObj.Status.State = ApplicationDeleting
		request.appObj.UpdatedAt = time.Now()
		// Block the command state in the queue and run the command to completion.
		queue.Started()
		tries := request.retryCount
		for tries >= 0 {
			_, err := chart.Remove()
			if err != nil {
				tries--
				glog.Errorf("remove retry count: %v", request.retryCount)
				glog.Errorf("remove failed on: %v", err)
				request.appObj.Status.State = ApplicationFailed
			} else {
				if *krak8sCfg.debug {
					glog.Infof("command execution success, tries: %d", tries)
				}
				tries = -1
				request.appObj.Status.State = ApplicationDeleted
				request.appObj.Status.DeployedAt = time.Now()
			}
			request.appObj.UpdatedAt = time.Now()
		}
		queue.Done()
	}

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

	glog.Infof("Queued task deleted: type: %s, name: %s, namespace: %s, queueing duration: %s, running duration %s",
		request.requestType.String(), request.projObj.Name, request.nsObj.Name, queue.QueuedDuration().String(), queue.RunningDuration().String())

	// ok to remove the request from the pending map
	r.mutex.Lock()
	delete(r.pendingRequests, index)
	r.mutex.Unlock()
	return
}

// ProjectRequest - submit project add request for processing.
func (r *Runner) ProjectRequest(action RequestType, ds *DataStore, proj *ProjectObject, ns *NamespaceObject, res *ResourceObject) RequestStatus {
	req := NewResourceRequest(action, ds, proj, ns, res)
	req.retryCount = 1
	queue.Submit(req.task)

	// add the request to the pending map
	r.mutex.Lock()
	r.index++
	r.pendingRequests[r.index] = req
	r.mutex.Unlock()
	r.sync <- r.index

	return Waiting
}

// ChartRequest - submit project add request for processing.
func (r *Runner) ChartRequest(action RequestType, ds *DataStore, proj *ProjectObject, ns *NamespaceObject, app *ApplicationObject) RequestStatus {
	req := NewChartRequest(action, ds, proj, ns, app)
	queue.Submit(req.task)

	// add the request to the pending map
	r.mutex.Lock()
	r.index++
	r.pendingRequests[r.index] = req
	r.mutex.Unlock()
	r.sync <- r.index

	return Waiting
}
