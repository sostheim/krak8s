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
	"sync"
	"time"
)

const (
	// OIDLength - Object ID (OID) string length
	OIDLength = 8
)

// API OjbectType Strings
const (
	// Project object type name
	Project = "project"
	// Namespace object type name
	Namespace = "namespace"
	// Resource object type name
	Resource = "Resource"
	// Application  object type name
	Application = "application"
)

// ObjectLink nested resource type
type ObjectLink struct {
	OID string
	URL string
}

// ProjectObject base resource type
type ProjectObject struct {
	OID        string
	ObjType    string
	Name       string
	CreatedAt  time.Time
	UpdatedAt  time.Time
	Namespaces []*ObjectLink
}

// NamespaceObject resource type
type NamespaceObject struct {
	OID          string
	ObjType      string
	Name         string
	CreatedAt    time.Time
	Resources    *ObjectLink
	Applications []*ObjectLink
}

// ApplicationStatusObject State strings
const (
	// Note: that the use UPPERCASE is intentional (it's a helm thing)

	// ApplicationUnknown state string
	ApplicationUnknown = "UNKNOWN"
	// ApplicationDeployed state string
	ApplicationDeployed = "DEPLOYED"
	// ApplicationDeleted state string
	ApplicationDeleted = "DELETED"
	// ApplicationSuperseded state string
	ApplicationSuperseded = "SUPERSEDED"
	// ApplicationDeleting state string
	ApplicationDeleting = "DELETING"
	// ApplicationFailed state string
	ApplicationFailed = "FAILED"
)

// ApplicationStatusObject nested object type
type ApplicationStatusObject struct {
	DeployedAt time.Time
	State      string
}

// ApplicationObject base resource type
type ApplicationObject struct {
	OID         string
	ObjType     string
	Name        string
	Version     string
	Config      string
	Registry    string
	CreatedAt   time.Time
	UpdatedAt   time.Time
	Status      *ApplicationStatusObject
	NamespaceID string
}

// ResourceObject State strings
const (
	// Note: that the use lowercase and understcroes is intentional

	// ResourceCreateRequested state string
	ResourceCreateRequested = "create_requested"
	// ResourceStarting state string
	ResourceStarting = "starting"
	// ResourceErrorStarting state string
	ResourceErrorStarting = "error_starting"
	// ResourceActive state string
	ResourceActive = "active"
	// ResourceDeleteRequested state string
	ResourceDeleteRequested = "delete_requested"
	// ResourceDeleting state string
	ResourceDeleting = "deleting"
	// ResourceErrorDeleting state string
	ResourceErrorDeleting = "error_deleting"
	// ResourceDeleted state string
	ResourceDeleted = "deleted"
)

// ResourceObject base resource type
type ResourceObject struct {
	OID          string
	ObjType      string
	NodePoolSize int
	CreatedAt    time.Time
	UpdatedAt    time.Time
	State        string
	NamespaceID  string
}

// DataStore in-memory data synchronization structure for API data.
type DataStore struct {
	sync.Mutex
	projects     map[string]*ProjectObject
	namespaces   map[string]*NamespaceObject
	resources    map[string]*ResourceObject
	applications map[string]*ApplicationObject
}

// NewDataStore initializes a new "DataStore"
func NewDataStore() *DataStore {
	ds := &DataStore{projects: nil, namespaces: nil, applications: nil, resources: nil}
	ds.Reset()
	return ds
}

// Reset removes all entries from the database. Mainly intended for tests.
func (ds *DataStore) Reset() {
	ds.projects = make(map[string]*ProjectObject)
	ds.namespaces = make(map[string]*NamespaceObject)
	ds.applications = make(map[string]*ApplicationObject)
	ds.resources = make(map[string]*ResourceObject)
}

// CheckedRandomHexString - generate a random hex string, and validate
// that the generated string is NOT a duplicate to a value already used.
func (ds *DataStore) CheckedRandomHexString() string {
	i := 9
	for i > 0 {
		tmp := RandHexString(OIDLength)
		if _, ok := ds.projects[tmp]; !ok {
			if _, ok := ds.namespaces[tmp]; !ok {
				if _, ok := ds.applications[tmp]; !ok {
					if _, ok := ds.resources[tmp]; !ok {
						return tmp
					}
				}
			}
		}
		i--
	}
	return ""
}

// NewProjectObject creates a default ProjectObject with a valid unique
// object id, type value, and created at timestamp.
func (ds *DataStore) NewProjectObject() *ProjectObject {
	obj := ProjectObject{
		OID:       ds.CheckedRandomHexString(),
		ObjType:   Project,
		CreatedAt: time.Now(),
	}
	if obj.OID == "" {
		return nil
	}
	ds.Lock()
	defer ds.Unlock()
	ds.projects[obj.OID] = &obj
	return &obj
}

// NewProject creates ProjectObject with given name and unique object id.
func (ds *DataStore) NewProject(name string) *ProjectObject {
	obj := ds.NewProjectObject()
	if obj == nil {
		return nil
	}
	obj.Name = name
	return obj
}

// ProjectsCollection returns all projects.
func (ds *DataStore) ProjectsCollection() []*ProjectObject {
	ds.Lock()
	defer ds.Unlock()
	i := 0
	collection := make([]*ProjectObject, len(ds.projects))
	for _, proj := range ds.projects {
		collection[i] = proj
		i++
	}
	return collection
}

// Project returns the rquested project object or, nil if not found.
func (ds *DataStore) Project(oid string) (*ProjectObject, bool) {
	ds.Lock()
	defer ds.Unlock()
	proj, ok := ds.projects[oid]
	return proj, ok
}

// AddProject add a project object to the data store.
func (ds *DataStore) AddProject(obj ProjectObject) {
	ds.Lock()
	defer ds.Unlock()
	ds.projects[obj.OID] = &obj
}

// DeleteProject removes the project and all subordinate objects.
func (ds *DataStore) DeleteProject(obj *ProjectObject) {
	for _, link := range ds.projects[obj.OID].Namespaces {
		ds.DeleteNamespace(ds.namespaces[link.OID])
	}

	ds.Lock()
	defer ds.Unlock()

	delete(ds.projects, obj.OID)
}

// NewNamespaceObject creates aa default NamespaceObject with a valid unique
// object id, type value and created at timestamp.
func (ds *DataStore) NewNamespaceObject() *NamespaceObject {
	obj := NamespaceObject{
		OID:          ds.CheckedRandomHexString(),
		ObjType:      Namespace,
		CreatedAt:    time.Now(),
		Resources:    nil,
		Applications: nil,
	}
	if obj.OID == "" {
		return nil
	}
	ds.Lock()
	defer ds.Unlock()
	ds.namespaces[obj.OID] = &obj
	return &obj
}

// NewNamespace creates NamespaceObject with given name
func (ds *DataStore) NewNamespace(name string) *NamespaceObject {
	obj := ds.NewNamespaceObject()
	if obj == nil {
		return nil
	}
	obj.Name = name
	return obj
}

// NamespacesCollection returns all Namespaces for a project
func (ds *DataStore) NamespacesCollection(projectOID string) []*NamespaceObject {
	proj, ok := ds.projects[projectOID]
	if !ok {
		return nil
	}

	ds.Lock()
	defer ds.Unlock()
	i := 0
	collection := make([]*NamespaceObject, len(proj.Namespaces))
	for _, link := range proj.Namespaces {
		collection[i] = ds.namespaces[link.OID]
		i++
	}
	return collection
}

// Namespace returns the rquested Namespace object or, nil if not found.
func (ds *DataStore) Namespace(oid string) (*NamespaceObject, bool) {
	ds.Lock()
	defer ds.Unlock()
	ns, ok := ds.namespaces[oid]
	return ns, ok
}

// AddNamespace add a Namespace object to the data store.
func (ds *DataStore) AddNamespace(obj NamespaceObject) {
	ds.Lock()
	defer ds.Unlock()
	ds.namespaces[obj.OID] = &obj
}

// DeleteNamespace removes the Namespace and all subordinate objects.
func (ds *DataStore) DeleteNamespace(obj *NamespaceObject) {
	for _, link := range ds.namespaces[obj.OID].Applications {
		ds.DeleteApplication(ds.applications[link.OID])
	}
	ds.DeleteResource(ds.namespaces[obj.OID].Resources.OID)

	ds.Lock()
	defer ds.Unlock()
	delete(ds.namespaces, obj.OID)
}

// NewApplicationObject creates aa default ApplicationObject with a valid
// unique object id. type value, and created at timestamp.
func (ds *DataStore) NewApplicationObject(nsOID string) *ApplicationObject {
	obj := ApplicationObject{
		OID:         ds.CheckedRandomHexString(),
		ObjType:     Application,
		CreatedAt:   time.Now(),
		Status:      &ApplicationStatusObject{State: ApplicationUnknown},
		NamespaceID: nsOID,
	}
	if obj.OID == "" {
		return nil
	}
	obj.UpdatedAt = obj.CreatedAt
	ds.Lock()
	defer ds.Unlock()
	ds.applications[obj.OID] = &obj
	return &obj
}

// NewApplication creates a new application resource.
func (ds *DataStore) NewApplication(namespace, name, version string, config, registry *string) *ApplicationObject {
	obj := ds.NewApplicationObject(namespace)
	if obj == nil {
		return nil
	}
	obj.Name = name
	obj.Version = version
	// Optional fields, may be unset (nil)
	if config != nil {
		obj.Config = *config
	}
	if registry != nil {
		obj.Registry = *registry
	}
	return obj
}

// Application returns the app with the given oid if found
func (ds *DataStore) Application(oid string) (*ApplicationObject, bool) {
	ds.Lock()
	defer ds.Unlock()
	app, ok := ds.applications[oid]
	return app, ok
}

// ApplicationsCollection return the collection of applications from the indicated namespace.
func (ds *DataStore) ApplicationsCollection(nsOID string) []*ApplicationObject {
	namespace, ok := ds.namespaces[nsOID]
	if !ok {
		return nil
	}
	i := 0
	collection := make([]*ApplicationObject, len(namespace.Applications))
	for _, link := range namespace.Applications {
		collection[i], ok = ds.Application(link.OID)
		if !ok {
			return nil
		}
		i++
	}
	return collection
}

// AddApplication add a Application object to the data store.
func (ds *DataStore) AddApplication(obj ApplicationObject) {
	ds.Lock()
	defer ds.Unlock()
	ds.applications[obj.OID] = &obj
}

// DeleteApplication deletes specified application
func (ds *DataStore) DeleteApplication(obj *ApplicationObject) {
	ds.Lock()
	defer ds.Unlock()
	delete(ds.applications, obj.OID)
}

// NewResourceObject creates a default ResourceObject with a valid unique id,
// type value, and created at timestamp.
func (ds *DataStore) NewResourceObject(nsOID string) *ResourceObject {
	obj := ResourceObject{
		OID:         ds.CheckedRandomHexString(),
		ObjType:     Resource,
		CreatedAt:   time.Now(),
		State:       ResourceCreateRequested,
		NamespaceID: nsOID,
	}
	if obj.OID == "" {
		return nil
	}
	ds.Lock()
	defer ds.Unlock()
	ds.resources[obj.OID] = &obj
	return &obj
}

// NewResource creates a new ResourceObject resource.
func (ds *DataStore) NewResource(namespace string, nodes int) *ResourceObject {
	obj := ds.NewResourceObject(namespace)
	if obj == nil {
		return nil
	}
	obj.NodePoolSize = nodes
	return obj
}

// Resource returns the app with the given oid if found
func (ds *DataStore) Resource(oid string) (*ResourceObject, bool) {
	ds.Lock()
	defer ds.Unlock()
	res, ok := ds.resources[oid]
	return res, ok
}

// ResourceObject return the resource object from the indicated namespace.
func (ds *DataStore) ResourceObject(nsOID string) (*ResourceObject, bool) {
	ns, ok := ds.namespaces[nsOID]
	if !ok {
		return nil, false
	}
	ds.Lock()
	defer ds.Unlock()
	return ds.Resource(ns.Resources.OID)
}

// AddResource add a rsource object to the data store.
func (ds *DataStore) AddResource(obj *ResourceObject) {
	ds.Lock()
	defer ds.Unlock()
	ds.resources[obj.OID] = obj
}

// DeleteResource deletes specified application
func (ds *DataStore) DeleteResource(resOID string) {
	ds.Lock()
	defer ds.Unlock()
	delete(ds.resources, resOID)
}
