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
	"k8s.io/client-go/kubernetes"

	"github.com/goadesign/goa"
	"github.com/goadesign/goa/middleware"

	"krak8s/app"
)

const (
	// APIVersion - URL path root for API versioning
	APIVersion = "/v1"
	// APIProjects - URL path segment for projects
	APIProjects = "/projects/"
	// APIApplications - URL path segment for applications
	APIApplications = "/applications/"
	// APICluster - URL path segment for cluster resources
	APICluster = "/cluster/"
	// APINamespaces - URL path segment for namespaces
	APINamespaces = "/namespaces/"
)

// API Server
type apiServer struct {
	cfg       *config
	clientset *kubernetes.Clientset
	server    *goa.Service
	ds        *DataStore
}

func newAPIServer(clientset *kubernetes.Clientset, cfg *config, backend *Runner) *apiServer {
	// Create and start the http router

	// create api server controller struct
	as := apiServer{
		clientset: clientset,
		cfg:       cfg,
		server:    goa.New("krak8s"),
		ds:        NewDataStore(),
	}

	// Mount middleware
	as.server.Use(middleware.RequestID())
	as.server.Use(middleware.LogRequest(true))
	as.server.Use(middleware.ErrorHandler(as.server, true))
	as.server.Use(middleware.Recover())

	// Create and Mount the resource controllers
	swagger := NewSwaggerController(as.server)
	app.MountSwaggerController(as.server, swagger)

	openapi := NewOpenapiController(as.server)
	app.MountOpenapiController(as.server, openapi)

	project := NewProjectController(as.server, as.ds)
	app.MountProjectController(as.server, project)

	ns := NewNamespaceController(as.server, as.ds)
	app.MountNamespaceController(as.server, ns)

	application := NewApplicationController(as.server, as.ds)
	app.MountApplicationController(as.server, application)

	cluster := NewClusterController(as.server, as.ds, backend)
	app.MountClusterController(as.server, cluster)

	health := NewHealthController(as.server)
	app.MountHealthController(as.server, health)

	return &as
}

func (as *apiServer) run() {
	if err := as.server.ListenAndServe(":8080"); err != nil {
		as.server.LogError("startup", "err", err)
	}
}
