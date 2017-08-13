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
	"krak8s/queue"
)

// Run back end tasks with queue semantics

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

// Request - task to run
type Request struct {
	task      *queue.Task
	request   RequestType
	name      string // either a project or chart name
	namespace string
	nodes     int
	version   string
	registry  string
	config    string
}

// NewRequest creates an request for processing
func NewRequest(req RequestType) *Request {
	return &Request{
		task:    queue.NewTask(),
		request: req,
	}
}

var pendingRequests []*Request

// AddProjectRequest - submit project add request for processing.
func AddProjectRequest(name, namespace string, nodes int) {
	req := NewRequest(AddProject)
	req.name = name
	req.namespace = namespace
	req.nodes = nodes
	queue.Submit(req.task)
	pendingRequests = append(pendingRequests, req)
}
