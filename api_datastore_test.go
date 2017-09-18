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
	"io/ioutil"
	"os"
	"testing"
)

var validDataStoreJSON = `
{
	"projects": {
	  "9b1c384c": {
		"oid": "9b1c384c",
		"objType": "project",
		"name": "saturn",
		"createdAt": "2017-09-13T22:19:26.947616505-07:00",
		"updatedAt": "0001-01-01T00:00:00Z",
		"namespaces": [{
		  "oid": "2f6356f0",
		  "url": "/v1/projects/9b1c384c/namespaces/2f6356f0"
		}]
	  },
	  "dbc5b124": {
		"oid": "dbc5b124",
		"objType": "project",
		"name": "neptune",
		"createdAt": "2017-09-13T22:14:55.708108344-07:00",
		"updatedAt": "0001-01-01T00:00:00Z",
		"namespaces": [{
		  "oid": "96f6162c",
		  "url": "/v1/projects/dbc5b124/namespaces/96f6162c"
		}, {
		  "oid": "c09db578",
		  "url": "/v1/projects/dbc5b124/namespaces/c09db578"
		}]
	  }
	},
	"namespaces": {
	  "2f6356f0": {
		"oid": "2f6356f0",
		"objType": "namespace",
		"name": "saturn-rings",
		"createdAt": "2017-09-13T22:19:51.394596951-07:00"
	  },
	  "96f6162c": {
		"oid": "96f6162c",
		"objType": "namespace",
		"name": "neptune-one",
		"createdAt": "2017-09-13T22:16:28.598556716-07:00",
		"resources": {
		  "oid": "4b3ff7db",
		  "url": "/v1/projects/dbc5b124/cluster/4b3ff7db"
		}
	  },
	  "c09db578": {
		"oid": "c09db578",
		"objType": "namespace",
		"name": "neptune-two",
		"createdAt": "2017-09-13T22:17:45.842641545-07:00",
		"resources": {
		  "oid": "6d9dd677",
		  "url": "/v1/projects/dbc5b124/cluster/6d9dd677"
		}
	  }
	},
	"resources": {
	  "4b3ff7db": {
		"oid": "4b3ff7db",
		"objType": "Resource",
		"nodePoolSize": 5,
		"createdAt": "2017-09-13T22:17:13.268069411-07:00",
		"updatedAt": "2017-09-13T22:17:13.271685248-07:00",
		"state": "active",
		"namespaceId": "96f6162c"
	  },
	  "6d9dd677": {
		"oid": "6d9dd677",
		"objType": "Resource",
		"nodePoolSize": 7,
		"createdAt": "2017-09-13T22:18:03.892229484-07:00",
		"updatedAt": "0001-01-01T00:00:00Z",
		"state": "create_requested",
		"namespaceId": "c09db578"
	  }
	}
  }
`

// missing comma delimiter between projects and namespaces
var invalidDataStoreJSON = `
{
	"projects": {
	  "9b1c384c": {
		"oid": "9b1c384c",
		"objType": "project",
		"name": "saturn",
		"createdAt": "2017-09-13T22:19:26.947616505-07:00",
		"updatedAt": "0001-01-01T00:00:00Z",
		"namespaces": [{
		  "oid": "2f6356f0",
		  "url": "/v1/projects/9b1c384c/namespaces/2f6356f0"
		}]
	  }
	}
	"namespaces": {
	  "2f6356f0": {
		"oid": "2f6356f0",
		"objType": "namespace",
		"name": "saturn-rings",
		"createdAt": "2017-09-13T22:19:51.394596951-07:00"
	  }
	}
  }
`

func TestNewDefaultDataStore(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test in short mode.")
	}

	ds := NewDataStore("")
	if ds == nil || ds.archive == nil || ds.persist != "" || ds.data.Applications == nil ||
		ds.data.Namespaces == nil || ds.data.Projects == nil || ds.data.Resources == nil {
		t.Error("NewDataStore(\"\") = nil, want: default reset object")
	}
	if len(ds.data.Projects) != 0 || len(ds.data.Namespaces) != 0 ||
		len(ds.data.Resources) != 0 || len(ds.data.Applications) != 0 {
		t.Errorf("NewDataStore(\"\") invalid dimensions P:%d/N:%d/R:%d/A:%d, want: P:0/N:0/R:0/A:0",
			len(ds.data.Projects), len(ds.data.Namespaces), len(ds.data.Resources), len(ds.data.Applications))
	}
}

func TestNewValidDataStoreLoad(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test in short mode.")
	}

	file, err := ioutil.TempFile(os.TempDir(), "test-vds")
	if err != nil {
		t.Error("TestNewValidDataStoreLoad() can't create temporary datastore file")
	}
	defer os.Remove(file.Name())

	if _, err := file.Write([]byte(validDataStoreJSON)); err != nil {
		t.Error("TestNewValidDataStoreLoad() can't write temporary datastore contents")
	}
	if err := file.Close(); err != nil {
		t.Error("TestNewValidDataStoreLoad() can't close temporary datastore file")
	}

	ds := NewDataStore(file.Name())
	if ds == nil || ds.archive == nil || ds.persist != file.Name() || ds.data.Applications == nil ||
		ds.data.Namespaces == nil || ds.data.Projects == nil || ds.data.Resources == nil {
		t.Errorf("NewDataStore(%s) = nil, want: valid datastore", file.Name())
	}
	if len(ds.data.Projects) != 2 || len(ds.data.Namespaces) != 3 ||
		len(ds.data.Resources) != 2 || len(ds.data.Applications) != 0 {
		t.Errorf("NewDataStore(\"\") invalid dimensions P:%d/N:%d/R:%d/A:%d, want: P:2/N:3/R:2/A:0",
			len(ds.data.Projects), len(ds.data.Namespaces), len(ds.data.Resources), len(ds.data.Applications))
	}
}

func TestNewInalidDataStoreLoad(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test in short mode.")
	}

	file, err := ioutil.TempFile(os.TempDir(), "test-invds")
	if err != nil {
		t.Error("TestNewValidDataStoreLoad() can't create temporary datastore file")
	}
	defer os.Remove(file.Name())

	if _, err := file.Write([]byte(invalidDataStoreJSON)); err != nil {
		t.Error("TestNewValidDataStoreLoad() can't write temporary datastore contents")
	}
	if err := file.Close(); err != nil {
		t.Error("TestNewValidDataStoreLoad() can't close temporary datastore file")
	}

	ds := NewDataStore(file.Name())
	if ds == nil || ds.archive == nil || ds.persist != file.Name() || ds.data.Applications == nil ||
		ds.data.Namespaces == nil || ds.data.Projects == nil || ds.data.Resources == nil {
		t.Errorf("NewDataStore(%s) = nil, want: valid datastore", file.Name())
	}
	if len(ds.data.Projects) != 0 || len(ds.data.Namespaces) != 0 ||
		len(ds.data.Resources) != 0 || len(ds.data.Applications) != 0 {
		t.Errorf("NewDataStore(\"\") invalid dimensions P:%d/N:%d/R:%d/A:%d, want: P:0/N:0/R:0/A:0",
			len(ds.data.Projects), len(ds.data.Namespaces), len(ds.data.Resources), len(ds.data.Applications))
	}
}

func TestNewProject(t *testing.T) {
	ds := NewDataStore("")
	if ds == nil {
		t.Errorf("NewDataStore() = nil, want: valid datastore")
	}
	obj := ds.NewProjectObject()
	if obj == nil {
		t.Errorf("NewProjectObject() = nil, want: valid project")
	}
	if obj.CreatedAt.IsZero() != false {
		t.Errorf("NewProjectObject() date = %v, want CratedAt != 0", obj.CreatedAt)
	}
	if obj.Name != "" {
		t.Errorf("NewProjectObject() Name = %s, want Name == ``", obj.Name)
	}
	if len(obj.Namespaces) != 0 {
		t.Errorf("NewProjectObject() namesapces = %d, want Namespaces = 0", len(obj.Namespaces))
	}
	if obj.ObjType != Project {
		t.Errorf("NewProjectObject() ObjType = %s, want: ObjType == Project", obj.ObjType)
	}
	if obj.OID == "" {
		t.Errorf("NewProjectObject(), OID = ``, want oid != ``")
	}
	if obj.UpdatedAt.IsZero() == false {
		t.Errorf("NewProjectObject()  update = %v, want UpdatedAt == Zero", obj.UpdatedAt)
	}
}

func TestNewNamedProject(t *testing.T) {
	ds := NewDataStore("")
	if ds == nil {
		t.Errorf("NewDataStore() = nil, want: valid datastore")
	}
	obj := ds.NewProject("test_object")
	if obj == nil {
		t.Errorf("NewProject() = nil, want: valid project")
	}
	if obj.CreatedAt.IsZero() == true {
		t.Error("NewProject() have CreatedAt == Zero, want CratedAt != Zero")
	}
	if obj.Name != "test_object" {
		t.Errorf("NewProject() Name = %s, want Name == `test_object`", obj.Name)
	}
	if len(obj.Namespaces) != 0 {
		t.Errorf("NewProject() namesapces = %d, want Namespaces = 0", len(obj.Namespaces))
	}
	if obj.ObjType != Project {
		t.Errorf("NewProject() ObjType = %s, want: ObjType == Project", obj.ObjType)
	}
	if obj.OID == "" {
		t.Errorf("NewProject(), OID = ``, want oid != ``")
	}
	if obj.UpdatedAt.IsZero() == true {
		t.Errorf("NewProject()  update = %v, want UpdatedAt == 0", obj.UpdatedAt)
	}
}

func TestProjectCollection(t *testing.T) {
	ds := NewDataStore("")
	if ds == nil {
		t.Errorf("NewDataStore() = nil, want: valid datastore")
	}
	go ds.Archiver()
	ds.NewProject("test_project_1")
	ds.NewProject("test_project_2")
	ds.NewProject("test_project_3")
	ds.NewProject("test_project_4")
	ds.NewProject("test_project_5")
	col := ds.ProjectsCollection()
	if len(col) != 5 {
		t.Errorf("ProjectsCollection() projects = %d, want projects = 5", len(col))
	}
}

func TestProject(t *testing.T) {
	ds := NewDataStore("")
	if ds == nil {
		t.Errorf("NewDataStore() = nil, want: valid datastore")
	}
	go ds.Archiver()
	ds.NewProject("test_project_1")
	ds.NewProject("test_project_2")
	obj := ds.NewProject("test_project_3")
	if obj == nil {
		t.Errorf("NewProject() = nil, want: valid project")
	}
	if obj.Name != "test_project_3" {
		t.Errorf("NewProject() Name = %s, want Name == `test_object`", obj.Name)
	}
	ds.NewProject("test_project_4")
	ds.NewProject("test_project_5")

	proj, found := ds.Project(obj.OID)
	if !found {
		t.Errorf("Project(%s) not found, want found", obj.OID)
	}
	if obj != proj {
		t.Errorf("Project(%s) obj != proj, want obj == proj", obj.OID)
	}

	_, found = ds.Project("badoid")
	if found {
		t.Errorf("Project(%s) found, want not found", "badoid")
	}
}

func TestNewNamespace(t *testing.T) {
	ds := NewDataStore("")
	if ds == nil {
		t.Errorf("NewDataStore() = nil, want: valid datastore")
	}
	ns := ds.NewNamespaceObject()
	if ns == nil {
		t.Errorf("NewNamespaceObject() = nil, want: valid Namespace")
	}
	if ns.ObjType != Namespace {
		t.Errorf("NewNamespaceObject(), ObjType = %s, want ObjType = %s", ns.ObjType, Namespace)
	}
	if ns.OID == "" {
		t.Errorf("NewNamespaceObject(), OID = ``, want oid != ``")
	}
	if ns.Name != "" {
		t.Errorf("NewNamespaceObject() Name = %s, want Name == ``", ns.Name)
	}
	if ns.CreatedAt.IsZero() != false {
		t.Errorf("NewNamespaceObject() date = %v, want CratedAt != 0", ns.CreatedAt)
	}
	if len(ns.Applications) != 0 {
		t.Errorf("NewNamespaceObject() applications = %d, want Applications = 0", len(ns.Applications))
	}
	if ns.Resources != nil {
		t.Error("NewNamespaceObject() resources exists, want Resources = nil")
	}
}

func TestNewNamedNamespace(t *testing.T) {
	ds := NewDataStore("")
	if ds == nil {
		t.Errorf("NewDataStore() = nil, want: valid datastore")
	}
	go ds.Archiver()
	obj := ds.NewProject("test_object")
	if obj == nil {
		t.Errorf("NewProject() = nil, want: valid project")
	}
	ns := ds.NewNamespace("test_namespace")
	if ns == nil {
		t.Errorf("NewNamespace() = nil, want: valid Namespace")
	}
	if ns.ObjType != Namespace {
		t.Errorf("NewNamespace(), ObjType = %s, want ObjType = %s", ns.ObjType, Namespace)
	}
	if ns.OID == "" {
		t.Errorf("NewNamespace(), OID = ``, want oid != ``")
	}
	if ns.Name != "test_namespace" {
		t.Errorf("NewNamespace() Name = %s, want Name == ``", ns.Name)
	}
	if ns.CreatedAt.IsZero() != false {
		t.Errorf("NewNamespace() date = %v, want CratedAt != 0", ns.CreatedAt)
	}
	if len(ns.Applications) != 0 {
		t.Errorf("NewNamespace() applications = %d, want Applications = 0", len(ns.Applications))
	}
	if ns.Resources != nil {
		t.Error("NewNamespace() resources exists, want Resources = nil")
	}
}

func TestNamespaceCollection(t *testing.T) {
	ds := NewDataStore("")
	if ds == nil {
		t.Errorf("NewDataStore() = nil, want: valid datastore")
	}
	go ds.Archiver()
	proj := ds.NewProject("test_project")
	if proj == nil {
		t.Errorf("NewProject() = nil, want: valid project")
	}
	for i := 0; i < 5; i++ {
		ns := ds.NewNamespace("test_namespace" + string(i))
		proj.Namespaces = append(proj.Namespaces, &ObjectLink{OID: ns.OID, URL: ""})
	}
	col := ds.NamespacesCollection(proj.OID)
	if len(col) != 5 {
		t.Errorf("NamespacesCollection() projects = %d, want namespaces = 5", len(col))
	}
}

func TestNamespace(t *testing.T) {
	ds := NewDataStore("")
	if ds == nil {
		t.Errorf("NewDataStore() = nil, want: valid datastore")
	}
	go ds.Archiver()
	proj := ds.NewProject("test_project")
	if proj == nil {
		t.Errorf("NewProject() = nil, want: valid project")
	}
	var obj *NamespaceObject
	for i := 0; i < 5; i++ {
		ns := ds.NewNamespace("test_namespace" + string(i))
		if i == 4 {
			obj = ns
		}
		proj.Namespaces = append(proj.Namespaces, &ObjectLink{OID: ns.OID, URL: ""})
	}
	nspc, found := ds.Namespace(obj.OID)
	if !found {
		t.Errorf("Namespace(%s) not found, want found", obj.OID)
	}
	if obj != nspc {
		t.Errorf("Namespace(%s) obj != nspc, want obj == nspc", obj.OID)
	}

	_, found = ds.Namespace("badoid")
	if found {
		t.Errorf("Namespace(%s) found, want not found", "badoid")
	}
}

func TestNewApplication(t *testing.T) {
	ds := NewDataStore("")
	if ds == nil {
		t.Errorf("NewDataStore() = nil, want: valid datastore")
	}
	ns := ds.NewNamespace("test_namespace")
	if ns == nil {
		t.Errorf("NewNamespace(%s), have: nil, want: valid Namespace", "test_namespace")
	}
	app := ds.NewApplicationObject(ns.OID)
	if app == nil {
		t.Errorf("NewApplicationObject(%s), have: nil, want: Application object", ns.OID)
	}
	if app.ObjType != Application {
		t.Errorf("NewApplicationObject(%s), have ObjType(%s), want: ObjType(%s)", ns.OID, app.ObjType, Application)
	}
	if app.OID == "" {
		t.Errorf("NewApplicationObject(%s), have OID = ``, want oid != ``", ns.OID)
	}
	if app.CreatedAt.IsZero() == true {
		t.Errorf("NewApplicationObject(%s) date = %v, want CratedAt != Zero", ns.OID, app.CreatedAt)
	}
	if app.NamespaceID != ns.OID {
		t.Errorf("NewApplicationObject(%s), have ns(%s), want ns(%s)", ns.OID, app.NamespaceID, ns.OID)
	}
}

func TestNewNamedApplication(t *testing.T) {
	ds := NewDataStore("")
	if ds == nil {
		t.Errorf("NewDataStore() = nil, want: valid datastore")
	}
	go ds.Archiver()
	obj := ds.NewProject("test_object")
	if obj == nil {
		t.Errorf("NewProject() = nil, want: valid project")
	}
	ns := ds.NewNamespace("test_namespace")
	if ns == nil {
		t.Errorf("NewNamespace(%s), have: nil, want: valid Namespace", "test_namespace")
	}
	chn := "test_channel"
	pwd := "test_password"
	app := ds.NewApplication(ns.OID, "test_deployment", "test_server", "test_registry",
		"test_chart", "test_version", &chn, nil, &pwd, nil, nil)
	if app == nil {
		t.Errorf("NewApplication(%s), have: nil, want: Application object", ns.OID)
	}
	if app.ObjType != Application {
		t.Errorf("NewApplication(%s), have ObjType(%s), want: ObjType(%s)", ns.OID, app.ObjType, Application)
	}
	if app.OID == "" {
		t.Errorf("NewApplication(%s), have OID = ``, want oid != ``", ns.OID)
	}
	if app.CreatedAt.IsZero() == true {
		t.Errorf("NewApplication(%s) date = %v, want CratedAt != Zero", ns.OID, app.CreatedAt)
	}
	if app.NamespaceID != ns.OID {
		t.Errorf("NewApplication(%s), have ns(%s), want ns(%s)", ns.OID, app.NamespaceID, ns.OID)
	}
	if app.Deployment != "test_deployment" {
		t.Errorf("NewApplication(%s), have deployment(%s), want deployment(`test_deployment`)", ns.OID, app.Deployment)
	}
	if app.Server != "test_server" {
		t.Errorf("NewApplication(%s), have server(%s), want server(`test_server`)", ns.OID, app.Server)
	}
	if app.ChartRegistry != "test_registry" {
		t.Errorf("NewApplication(%s), have registry(%s), want registry(`test_registry`)", ns.OID, app.ChartRegistry)
	}
	if app.ChartName != "test_chart" {
		t.Errorf("NewApplication(%s), have chart(%s), want chart(`test_chart`)", ns.OID, app.ChartName)
	}
	if app.ChartVersion != "test_version" {
		t.Errorf("NewApplication(%s), have version(%s), want version(`test_version`)", ns.OID, app.ChartVersion)
	}
	if app.Channel != "test_channel" {
		t.Errorf("NewApplication(%s), have channel(%s), want channel(`test_channel`)", ns.OID, app.Channel)
	}
	if app.Username != "" {
		t.Errorf("NewApplication(%s), have username(%s), want username(``)", ns.OID, app.Username)
	}
	if app.Password != "test_password" {
		t.Errorf("NewApplication(%s), have password(%s), want password(`test_password`)", ns.OID, app.Password)
	}
	if app.Config != "" {
		t.Errorf("NewApplication(%s), have config(%s), want config(``)", ns.OID, app.Config)
	}
	if app.JSONValues != "" {
		t.Errorf("NewApplication(%s), have json(%s), want json(``)", ns.OID, app.JSONValues)
	}
}

func TestApplicationCollection(t *testing.T) {
	ds := NewDataStore("")
	if ds == nil {
		t.Errorf("NewDataStore() = nil, want: valid datastore")
	}
	go ds.Archiver()
	proj := ds.NewProject("test_project")
	if proj == nil {
		t.Errorf("NewProject(), have: nil, want: project object")
	}
	ns := ds.NewNamespace("test_namespace")
	if ns == nil {
		t.Errorf("NewNamespace(%s), have: nil, want: namespace object", "test_namespace")
	}
	for i := 0; i < 5; i++ {
		app := ds.NewApplication(ns.OID, "test_deployment", "test_server", "test_registry",
			"test_chart", "test_version", nil, nil, nil, nil, nil)
		ns.Applications = append(ns.Applications, &ObjectLink{OID: app.OID, URL: ""})
	}
	col := ds.ApplicationsCollection(ns.OID)
	if len(col) != 5 {
		t.Errorf("ApplicationsCollection(%s) len(applications) = %d, want len(applications) = 5", ns.OID, len(col))
	}
}

func TestApplication(t *testing.T) {
	ds := NewDataStore("")
	if ds == nil {
		t.Errorf("NewDataStore() = nil, want: valid datastore")
	}
	go ds.Archiver()
	proj := ds.NewProject("test_project")
	if proj == nil {
		t.Errorf("NewProject(), have: nil, want: project object")
	}
	ns := ds.NewNamespace("test_namespace")
	if ns == nil {
		t.Errorf("NewNamespace(%s), have: nil, want: namespace object", "test_namespace")
	}
	var obj *ApplicationObject
	for i := 0; i < 5; i++ {
		app := ds.NewApplication(ns.OID, "test_deployment", "test_server", "test_registry",
			"test_chart", "test_version", nil, nil, nil, nil, nil)
		if i == 2 {
			obj = app
		}
		ns.Applications = append(ns.Applications, &ObjectLink{OID: app.OID, URL: ""})
	}
	app, found := ds.Application(obj.OID)
	if !found {
		t.Errorf("Application(%s) not found, want found", obj.OID)
	}
	if obj != app {
		t.Errorf("Application(%s) obj != app, want obj == app", obj.OID)
	}

	_, found = ds.Application("badoid")
	if found {
		t.Errorf("Application(%s) found, want not found", "badoid")
	}
}

func TestNewResource(t *testing.T) {
	ds := NewDataStore("")
	if ds == nil {
		t.Errorf("NewDataStore() = nil, want: valid datastore")
	}
	ns := ds.NewNamespace("test_namespace")
	if ns == nil {
		t.Errorf("NewNamespace(%s), have: nil, want: valid Namespace", "test_namespace")
	}
	res := ds.NewResourceObject(ns.OID)
	if res == nil {
		t.Errorf("NewResourceObject(%s), have: nil, want: Resource object", ns.OID)
	}
	if res.ObjType != Resource {
		t.Errorf("NewResourceObject(%s), have ObjType(%s), want: ObjType(%s)", ns.OID, res.ObjType, Resource)
	}
	if res.OID == "" {
		t.Errorf("NewResourceObject(%s), have OID = ``, want oid != ``", ns.OID)
	}
	if res.CreatedAt.IsZero() == true {
		t.Errorf("NewResourceObject(%s) date = %v, want CratedAt != Zero", ns.OID, res.CreatedAt)
	}
	if res.NamespaceID != ns.OID {
		t.Errorf("NewResourceObject(%s), have ns(%s), want ns(%s)", ns.OID, res.NamespaceID, ns.OID)
	}
}

func TestNewNamedResource(t *testing.T) {
	ds := NewDataStore("")
	if ds == nil {
		t.Errorf("NewDataStore() = nil, want: valid datastore")
	}
	go ds.Archiver()
	obj := ds.NewProject("test_object")
	if obj == nil {
		t.Errorf("NewProject() = nil, want: valid project")
	}
	ns := ds.NewNamespace("test_namespace")
	if ns == nil {
		t.Errorf("NewNamespace(%s), have: nil, want: valid Namespace", "test_namespace")
	}
	res := ds.NewResource(ns.OID, 3)
	if res == nil {
		t.Errorf("NewResource(%s), have: nil, want: Resource object", ns.OID)
	}
	if res.ObjType != Resource {
		t.Errorf("NewResource(%s), have ObjType(%s), want ObjType(%s)", ns.OID, res.ObjType, Resource)
	}
	if res.OID == "" {
		t.Errorf("NewResource(%s), have OID = ``, want oid != ``", ns.OID)
	}
	if res.CreatedAt.IsZero() == true {
		t.Errorf("NewResource(%s) date = %v, want CratedAt != Zero", ns.OID, res.CreatedAt)
	}
	if res.NamespaceID != ns.OID {
		t.Errorf("NewResource(%s), have ns oid(%s), want ns oid(%s)", ns.OID, res.NamespaceID, ns.OID)
	}
	if res.NodePoolSize != 3 {
		t.Errorf("NewResource(%s), have node pool size(%d), want node pool size(3)", ns.OID, res.NodePoolSize)
	}
}

func TestResource(t *testing.T) {
	ds := NewDataStore("")
	if ds == nil {
		t.Errorf("NewDataStore() = nil, want: valid datastore")
	}
	go ds.Archiver()
	proj := ds.NewProject("test_project")
	if proj == nil {
		t.Errorf("NewProject(), have: nil, want: project object")
	}
	ns := ds.NewNamespace("test_namespace")
	if ns == nil {
		t.Errorf("NewNamespace(%s), have: nil, want: namespace object", "test_namespace")
	}
	res := ds.NewResource(ns.OID, 7)
	if res == nil {
		t.Errorf("NewResource(%s), have: nil, want: Resource object", ns.OID)
	}
	ns.Resources = &ObjectLink{OID: res.OID, URL: ""}

	rsrc, found := ds.Resource(res.OID)
	if !found {
		t.Errorf("Resource(%s) not found, want found", res.OID)
	}
	if res != rsrc {
		t.Errorf("Resource(%s) res != rsrc, want res == rsrc", res.OID)
	}

	_, found = ds.Resource("badoid")
	if found {
		t.Errorf("Resource(%s) found, want not found", "badoid")
	}
}

func TestOIDGenerator(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test in short mode.")
	}
	// Generate 50M Unique Id's without a duplicate (takes ~60s)
	iterations := 50000000
	ds := NewDataStore("")
	if ds == nil {
		t.Errorf("NewDataStore() = nil, want: valid datastore")
	}
	go ds.Archiver()
	for i := 0; i < iterations; i++ {
		obj := ds.NewProjectObject()
		if obj == nil {
			t.Errorf("NewProjectObject() = nil, want: valid project")
		}
	}
	col := ds.ProjectsCollection()
	if len(col) != iterations {
		t.Errorf("ProjectsCollection() projects = %d, want projects = %d", len(col), iterations)
	}
}

func BenchmarkOIDGenerator(b *testing.B) {
	// Generate 10M Unique Id's without a duplicate
	ds := NewDataStore("")
	if ds == nil {
		b.Errorf("NewDataStore() = nil, want: valid datastore")
	}
	go ds.Archiver()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		proj := ds.NewProjectObject()
		if proj == nil {
			b.Error("NewProjectObject(), have nil, want Project object")
		}
		ns := ds.NewNamespaceObject()
		if ns == nil {
			b.Error("NewNamespaceObject(), have nil, want Namespace object")
			continue
		}
		res := ds.NewResourceObject(ns.OID)
		if res == nil {
			b.Errorf("NewResourceObject(%s), have nil, want Resource object", ns.OID)
		}
		app := ds.NewApplicationObject(ns.OID)
		if app == nil {
			b.Errorf("NewApplicationObject(%s), have nil, want Application object", ns.OID)
		}
	}
}
