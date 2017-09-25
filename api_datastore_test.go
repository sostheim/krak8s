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

const neptuneProjectID = "dbc5b124"
const validProjects = 2
const validNamesapces = 3
const validResources = 2
const validApps = 0

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

const illinoisProjectID = "6f6c206e"
const urbanaNamespaceID = "fc0c6a6b"
const urbanaMongoAppID = "7412131c"
const urbanaClusterID = "3e635557"
const geoProjects = 5
const geoNamesapces = 5
const geoResources = 3
const geoApps = 4

var geographJSON = `
{
	"projects": {
		"17bc9fa1": {
			"oid": "17bc9fa1",
			"objType": "project",
			"name": "texas",
			"createdAt": "2017-09-18T23:35:37.855666965-07:00",
			"updatedAt": "2017-09-18T23:35:37.855667561-07:00",
			"namespaces": [{
				"oid": "3a268717",
				"url": "/v1/projects/17bc9fa1/namespaces/3a268717"
			}]
		},
		"6f6c206e": {
			"oid": "6f6c206e",
			"objType": "project",
			"name": "illinois",
			"createdAt": "2017-09-18T22:17:45.866844928-07:00",
			"updatedAt": "2017-09-18T22:17:45.866845523-07:00",
			"namespaces": [{
				"oid": "fc0c6a6b",
				"url": "/v1/projects/6f6c206e/namespaces/fc0c6a6b"
			}, {
				"oid": "9b4ed8d3",
				"url": "/v1/projects/6f6c206e/namespaces/9b4ed8d3"
			}]
		},
		"72e3b7bd": {
			"oid": "72e3b7bd",
			"objType": "project",
			"name": "utah",
			"createdAt": "2017-09-20T00:27:51.486993229-07:00",
			"updatedAt": "2017-09-20T00:27:51.486993753-07:00"
		},
		"a397fe98": {
			"oid": "a397fe98",
			"objType": "project",
			"name": "alaska",
			"createdAt": "2017-09-18T21:45:52.039344512-07:00",
			"updatedAt": "2017-09-18T21:45:52.039346014-07:00",
			"namespaces": [{
				"oid": "c37c08c0",
				"url": "/v1/projects/a397fe98/namespaces/c37c08c0"
			}]
		},
		"d16c6120": {
			"oid": "d16c6120",
			"objType": "project",
			"name": "florida",
			"createdAt": "2017-09-20T00:29:11.154709289-07:00",
			"updatedAt": "2017-09-20T00:29:11.154710127-07:00",
			"namespaces": [{
				"oid": "7041c0fd",
				"url": "/v1/projects/d16c6120/namespaces/7041c0fd"
			}]
		}
	},
	"namespaces": {
		"3a268717": {
			"oid": "3a268717",
			"objType": "namespace",
			"name": "texas-houston",
			"createdAt": "2017-09-18T23:36:27.401279441-07:00"
		},
		"7041c0fd": {
			"oid": "7041c0fd",
			"objType": "namespace",
			"name": "florida-panhandle",
			"createdAt": "2017-09-20T00:29:43.069764356-07:00"
		},
		"9b4ed8d3": {
			"oid": "9b4ed8d3",
			"objType": "namespace",
			"name": "illinois-chicago",
			"createdAt": "2017-09-19T23:30:50.615063003-07:00",
			"resources": {
				"oid": "66b9fc13",
				"url": "/v1/projects/6f6c206e/cluster/66b9fc13"
			},
			"applications": [{
				"oid": "58498d77",
				"url": "/v1/projects/6f6c206e/applications/58498d77"
			}, {
				"oid": "00acc901",
				"url": "/v1/projects/6f6c206e/applications/00acc901"
			}]
		},
		"c37c08c0": {
			"oid": "c37c08c0",
			"objType": "namespace",
			"name": "alaska-juneau",
			"createdAt": "2017-09-18T21:46:10.981700708-07:00",
			"resources": {
				"oid": "f3c4c520",
				"url": "/v1/projects/a397fe98/cluster/f3c4c520"
			},
			"applications": [{
				"oid": "33f77ac2",
				"url": "/v1/projects/a397fe98/applications/33f77ac2"
			}]
		},
		"fc0c6a6b": {
			"oid": "fc0c6a6b",
			"objType": "namespace",
			"name": "illinois-urbana",
			"createdAt": "2017-09-18T22:18:13.432763982-07:00",
			"resources": {
				"oid": "3e635557",
				"url": "/v1/projects/6f6c206e/cluster/3e635557"
			},
			"applications": [{
				"oid": "7412131c",
				"url": "/v1/projects/6f6c206e/applications/7412131c"
			}]
		}
	},
	"resources": {
		"3e635557": {
			"oid": "3e635557",
			"objType": "Resource",
			"nodePoolSize": 5,
			"createdAt": "2017-09-18T22:18:49.396831343-07:00",
			"updatedAt": "2017-09-18T22:18:49.397778692-07:00",
			"state": "active",
			"namespaceId": "fc0c6a6b"
		},
		"66b9fc13": {
			"oid": "66b9fc13",
			"objType": "Resource",
			"nodePoolSize": 7,
			"createdAt": "2017-09-19T23:33:07.535247-07:00",
			"updatedAt": "0001-01-01T00:00:00Z",
			"state": "create_requested",
			"namespaceId": "9b4ed8d3"
		},
		"f3c4c520": {
			"oid": "f3c4c520",
			"objType": "Resource",
			"nodePoolSize": 5,
			"createdAt": "2017-09-18T21:46:46.573235148-07:00",
			"updatedAt": "2017-09-18T21:46:46.576602106-07:00",
			"state": "active",
			"namespaceId": "c37c08c0"
		}
	},
	"applications": {
		"00acc901": {
			"oid": "00acc901",
			"objType": "application",
			"namespaceId": "9b4ed8d3",
			"deploymentName": "illinois-chicago-mongodb",
			"resgistryServer": "quay.io",
			"chartRegistry": "samsung_cnct",
			"chartName": "mongodb-replicaset",
			"chartVersion": "latest",
			"channel": "stable",
			"createdAt": "2017-09-19T23:51:41.482334361-07:00",
			"updatedAt": "2017-09-19T23:51:41.484729279-07:00",
			"status": {
				"deployedAt": "2017-09-19T23:51:41.484729262-07:00",
				"state": "DEPLOYED"
			}
		},
		"33f77ac2": {
			"oid": "33f77ac2",
			"objType": "application",
			"namespaceId": "c37c08c0",
			"deploymentName": "alaska-juneau-mongodb",
			"resgistryServer": "quay.io",
			"chartRegistry": "samsung_cnct",
			"chartName": "mongodb-replicaset",
			"chartVersion": "latest",
			"channel": "stable",
			"createdAt": "2017-09-18T21:56:03.633589716-07:00",
			"updatedAt": "2017-09-18T21:56:03.635134882-07:00",
			"status": {
				"deployedAt": "2017-09-18T21:56:03.635134774-07:00",
				"state": "DEPLOYED"
			}
		},
		"58498d77": {
			"oid": "58498d77",
			"objType": "application",
			"namespaceId": "9b4ed8d3",
			"deploymentName": "illinois-chicago-web",
			"resgistryServer": "quay.io",
			"chartRegistry": "samsung_cnct",
			"chartName": "apache",
			"chartVersion": "latest",
			"channel": "stable",
			"createdAt": "2017-09-19T23:49:33.757606234-07:00",
			"updatedAt": "2017-09-19T23:49:33.760895338-07:00",
			"status": {
				"deployedAt": "2017-09-19T23:49:33.760895017-07:00",
				"state": "DEPLOYED"
			}
		},
		"7412131c": {
			"oid": "7412131c",
			"objType": "application",
			"namespaceId": "fc0c6a6b",
			"deploymentName": "illinois-urbana-mongodb",
			"resgistryServer": "quay.io",
			"chartRegistry": "samsung_cnct",
			"chartName": "mongodb-replicaset",
			"chartVersion": "latest",
			"channel": "stable",
			"createdAt": "2017-09-18T22:32:31.325098155-07:00",
			"updatedAt": "2017-09-18T22:32:31.325669175-07:00",
			"status": {
				"deployedAt": "2017-09-18T22:32:31.325669159-07:00",
				"state": "DEPLOYED"
			}
		}
	}
}
`

func TestNewDefaultDataStore(t *testing.T) {
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
	if len(ds.data.Projects) != validProjects || len(ds.data.Namespaces) != validNamesapces ||
		len(ds.data.Resources) != validResources || len(ds.data.Applications) != validApps {
		t.Errorf("NewDataStore(\"\") invalid dimensions P:%d/N:%d/R:%d/A:%d, want: P:%d/N:%d/R:%d/A:%d",
			len(ds.data.Projects), len(ds.data.Namespaces), len(ds.data.Resources), len(ds.data.Applications),
			validProjects, validNamesapces, validResources, validApps)
	}
}

func TestNewInalidDataStoreLoad(t *testing.T) {
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
	defer close(ds.archive)
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
	defer close(ds.archive)
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

func TestDeleteProjectObject(t *testing.T) {
	ds := NewDataStore("")
	if ds == nil {
		t.Errorf("NewDataStore() have nil, want DataStore object")
	}
	go ds.Archiver()
	defer close(ds.archive)
	ds.NewProjectObject()
	ds.NewProjectObject()
	obj := ds.NewProjectObject()
	if obj == nil {
		t.Errorf("NewProjectObject() have nil, want Project object")
	}
	ds.NewProjectObject()
	ds.NewProjectObject()
	pre := len(ds.data.Projects)
	ds.DeleteProject(obj)
	if len(ds.data.Projects) != pre-1 {
		t.Errorf("TestDeleteProjectObjects() have len(%d), want len(%d)",
			len(ds.data.Projects), pre-1)
	}
}

func TestDeleteProject(t *testing.T) {
	file, err := ioutil.TempFile(os.TempDir(), "test-dp")
	if err != nil {
		t.Errorf("TestDeleteProject() have err: %v, want valid file", err)
	}
	defer os.Remove(file.Name())

	if _, err := file.Write([]byte(validDataStoreJSON)); err != nil {
		t.Errorf("TestDeleteProject() write temporary datastore err: %v", err)
	}
	if err := file.Close(); err != nil {
		t.Errorf("TestDeleteProject() close temporary datastore file err: %v", err)
	}

	ds := NewDataStore(file.Name())
	go ds.Archiver()
	defer close(ds.archive)
	preproj, prens, prers, preapp := len(ds.data.Projects), len(ds.data.Namespaces), len(ds.data.Resources), len(ds.data.Applications)
	proj, found := ds.Project(neptuneProjectID)
	if proj == nil || !found {
		t.Errorf("Project() have nil/not found, want Project object")
	}
	ds.DeleteProject(proj)
	if len(ds.data.Projects) != preproj-1 {
		t.Errorf("TestDeleteProject() have proj len(%d), want proj len(%d)",
			len(ds.data.Projects), preproj-1)
	}
	// project has 2 namespaces
	if len(ds.data.Namespaces) != prens-2 {
		t.Errorf("TestDeleteProject() have ns len(%d), want ns len(%d)",
			len(ds.data.Namespaces), prens-2)
	}
	// each namespacce has a resource
	if len(ds.data.Resources) != prers-2 {
		t.Errorf("TestDeleteProject() have resource len(%d), want resource len(%d)",
			len(ds.data.Resources), prers-2)
	}
	if len(ds.data.Applications) != preapp {
		t.Errorf("TestDeleteProject() have applications len(%d), want applications len(%d)",
			len(ds.data.Resources), preapp)
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
	defer close(ds.archive)
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
	defer close(ds.archive)
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
	defer close(ds.archive)
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

func TestDeleteNamespaceObject(t *testing.T) {
	ds := NewDataStore("")
	if ds == nil {
		t.Errorf("NewDataStore() have nil, want DataStore object")
	}
	go ds.Archiver()
	defer close(ds.archive)
	obj := ds.NewNamespaceObject()
	if obj == nil {
		t.Errorf("NewNamespaceObject() have nil, want Namespace object")
	}
	ds.NewNamespaceObject()
	ds.NewNamespaceObject()
	ds.NewNamespaceObject()
	ds.NewNamespaceObject()
	pre := len(ds.data.Namespaces)
	ds.DeleteNamespace(obj)
	if len(ds.data.Namespaces) != pre-1 {
		t.Errorf("TestDeleteNamespaceObject() have len(%d), want len(%d)",
			len(ds.data.Namespaces), pre-1)
	}
}

func TestDeleteNamespace(t *testing.T) {
	file, err := ioutil.TempFile(os.TempDir(), "test-dns")
	if err != nil {
		t.Errorf("TestDeleteNamespace() have err: %v, want valid file", err)
	}
	defer os.Remove(file.Name())

	if _, err := file.Write([]byte(geographJSON)); err != nil {
		t.Errorf("TestDeleteNamespace() write temporary datastore err: %v", err)
	}
	if err := file.Close(); err != nil {
		t.Errorf("TestDeleteNamespace() close temporary datastore file err: %v", err)
	}

	ds := NewDataStore(file.Name())
	go ds.Archiver()
	defer close(ds.archive)
	preproj, prens, prers, preapp := len(ds.data.Projects), len(ds.data.Namespaces), len(ds.data.Resources), len(ds.data.Applications)

	ns, found := ds.Namespace(urbanaNamespaceID)
	if ns == nil || !found {
		t.Errorf("Namespace() have nil/not found, want Project object")
	}
	ds.DeleteNamespace(ns)
	// project collection remains unchanged
	if len(ds.data.Projects) != preproj {
		t.Errorf("TestDeleteNamespace() have proj len(%d), want proj len(%d)",
			len(ds.data.Projects), preproj-1)
	}
	if len(ds.data.Namespaces) != prens-1 {
		t.Errorf("TestDeleteNamespace() have ns len(%d), want ns len(%d)",
			len(ds.data.Namespaces), prens-1)
	}
	if len(ds.data.Resources) != prers-1 {
		t.Errorf("TestDeleteNamespace() have resource len(%d), want resource len(%d)",
			len(ds.data.Resources), prers-1)
	}
	if len(ds.data.Applications) != preapp-1 {
		t.Errorf("TestDeleteNamespace() have applications len(%d), want applications len(%d)",
			len(ds.data.Resources), preapp-1)
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
	defer close(ds.archive)
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
	defer close(ds.archive)
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
	defer close(ds.archive)
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

func TestDeleteApplicationObject(t *testing.T) {
	ds := NewDataStore("")
	if ds == nil {
		t.Errorf("NewDataStore() have nil, want DataStore object")
	}
	go ds.Archiver()
	defer close(ds.archive)
	ns := ds.NewNamespaceObject()
	if ns == nil {
		t.Errorf("NewNamespaceObject() have nil, want Namespace object")
	}
	ds.NewApplicationObject(ns.OID)
	ds.NewApplicationObject(ns.OID)
	ds.NewApplicationObject(ns.OID)
	ds.NewApplicationObject(ns.OID)
	app := ds.NewApplicationObject(ns.OID)
	if app == nil {
		t.Errorf("NewApplicationObject(%s), have nil, want Application object", ns.OID)
	}
	pre := len(ds.data.Applications)
	ds.DeleteApplication(app)
	if len(ds.data.Applications) != pre-1 {
		t.Errorf("TestDeleteApplicationObject() have len(%d), want len(%d)",
			len(ds.data.Applications), pre-1)
	}
}

func TestDeleteApplication(t *testing.T) {
	file, err := ioutil.TempFile(os.TempDir(), "test-da")
	if err != nil {
		t.Errorf("TestDeleteApplication() have err: %v, want valid file", err)
	}
	defer os.Remove(file.Name())

	if _, err := file.Write([]byte(geographJSON)); err != nil {
		t.Errorf("TestDeleteApplication() write temporary datastore err: %v", err)
	}
	if err := file.Close(); err != nil {
		t.Errorf("TestDeleteApplication() close temporary datastore file err: %v", err)
	}

	ds := NewDataStore(file.Name())
	go ds.Archiver()
	defer close(ds.archive)
	preproj, prens, prers, preapp := len(ds.data.Projects), len(ds.data.Namespaces), len(ds.data.Resources), len(ds.data.Applications)

	app, found := ds.Application(urbanaMongoAppID)
	if app == nil || !found {
		t.Errorf("Application() have nil/not found, want Application object")
	}
	ds.DeleteApplication(app)
	// project, namespace, resources collection remains unchanged
	if len(ds.data.Projects) != preproj {
		t.Errorf("TestDeleteApplication() have proj len(%d), want proj len(%d)",
			len(ds.data.Projects), preproj-1)
	}
	if len(ds.data.Namespaces) != prens {
		t.Errorf("TestDeleteApplication() have ns len(%d), want ns len(%d)",
			len(ds.data.Namespaces), prens)
	}
	if len(ds.data.Resources) != prers {
		t.Errorf("TestDeleteApplication() have resource len(%d), want resource len(%d)",
			len(ds.data.Resources), prers)
	}
	if len(ds.data.Applications) != preapp-1 {
		t.Errorf("TestDeleteApplication() have applications len(%d), want applications len(%d)",
			len(ds.data.Resources), preapp-1)
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
	defer close(ds.archive)
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
	defer close(ds.archive)
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

func TestDeleteResourceObject(t *testing.T) {
	ds := NewDataStore("")
	if ds == nil {
		t.Errorf("NewDataStore() have nil, want DataStore object")
	}
	go ds.Archiver()
	defer close(ds.archive)
	ns := ds.NewNamespaceObject()
	if ns == nil {
		t.Errorf("NewNamespaceObject() have nil, want Namespace object")
	}
	res := ds.NewResourceObject(ns.OID)
	if res == nil {
		t.Errorf("NewResourceObject(%s), have nil, want Resource object", ns.OID)
	}
	pre := len(ds.data.Resources)
	ds.DeleteResource(res)
	if len(ds.data.Resources) != pre-1 {
		t.Errorf("TestDeleteResourceObject() have len(%d), want len(%d)",
			len(ds.data.Resources), pre-1)
	}
}

func TestDeleteResource(t *testing.T) {
	file, err := ioutil.TempFile(os.TempDir(), "test-da")
	if err != nil {
		t.Errorf("TestDeleteResource() have err: %v, want valid file", err)
	}
	defer os.Remove(file.Name())

	if _, err := file.Write([]byte(geographJSON)); err != nil {
		t.Errorf("TestDeleteResource() write temporary datastore err: %v", err)
	}
	if err := file.Close(); err != nil {
		t.Errorf("TestDeleteResource() close temporary datastore file err: %v", err)
	}

	ds := NewDataStore(file.Name())
	go ds.Archiver()
	defer close(ds.archive)
	preproj, prens, prers, preapp := len(ds.data.Projects), len(ds.data.Namespaces), len(ds.data.Resources), len(ds.data.Applications)

	res, found := ds.Resource(urbanaClusterID)
	if res == nil || !found {
		t.Errorf("Resource() have nil/not found, want Resource object")
	}
	ds.DeleteResource(res)
	// project, namespace, application collections remain unchanged
	if len(ds.data.Projects) != preproj {
		t.Errorf("TestDeleteResource() have proj len(%d), want proj len(%d)",
			len(ds.data.Projects), preproj-1)
	}
	if len(ds.data.Namespaces) != prens {
		t.Errorf("TestDeleteResource() have ns len(%d), want ns len(%d)",
			len(ds.data.Namespaces), prens)
	}
	if len(ds.data.Applications) != preapp {
		t.Errorf("TestDeleteResource() have applications len(%d), want applications len(%d)",
			len(ds.data.Resources), preapp)
	}
	if len(ds.data.Resources) != prers-1 {
		t.Errorf("TestDeleteApplication() have resource len(%d), want resource len(%d)",
			len(ds.data.Resources), prers-1)
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
	defer close(ds.archive)
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
	defer close(ds.archive)
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
