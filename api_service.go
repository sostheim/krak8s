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
	"time"

	"k8s.io/client-go/kubernetes"

	"github.com/goadesign/goa"
	"github.com/goadesign/goa/middleware"

	"krak8s/app"
)

var (
	resyncPeriod = 30 * time.Second
)

// API Server
type apiServer struct {
	// Command line / environment supplied configuration values
	cfg       *config
	clientset *kubernetes.Clientset
	server    *goa.Service
}

func newAPIServer(clientset *kubernetes.Clientset, cfg *config) *apiServer {
	// Create and start the http router

	// create api server controller struct
	as := apiServer{
		clientset: clientset,
		cfg:       cfg,
		server:    goa.New("krak8s"),
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

	project := NewProjectController(as.server)
	app.MountProjectController(as.server, project)

	ns := NewNamespaceController(as.server)
	app.MountNamespaceController(as.server, ns)

	mongo := NewMongoController(as.server)
	app.MountMongoController(as.server, mongo)

	chart := NewChartController(as.server)
	app.MountChartController(as.server, chart)

	cluster := NewClusterController(as.server)
	app.MountClusterController(as.server, cluster)

	return &as
}

func (as *apiServer) run() {
	if err := as.server.ListenAndServe(":8080"); err != nil {
		as.server.LogError("startup", "err", err)
	}
}
