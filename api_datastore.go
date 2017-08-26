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
	"encoding/json"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"sync"
	"time"

	"github.com/golang/glog"
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
	OID string `json:"oid,omitempty"`
	URL string `json:"url,omitempty"`
}

// ProjectObject base resource type
type ProjectObject struct {
	OID        string        `json:"oid,omitempty"`
	ObjType    string        `json:"objType,omitempty"`
	Name       string        `json:"name,omitempty"`
	CreatedAt  time.Time     `json:"createdAt,omitempty"`
	UpdatedAt  time.Time     `json:"updatedAt,omitempty"`
	Namespaces []*ObjectLink `json:"namespaces,omitempty"`
}

// NamespaceObject resource type
type NamespaceObject struct {
	OID          string        `json:"oid,omitempty"`
	ObjType      string        `json:"objType,omitempty"`
	Name         string        `json:"name,omitempty"`
	CreatedAt    time.Time     `json:"createdAt,omitempty"`
	Resources    *ObjectLink   `json:"resources,omitempty"`
	Applications []*ObjectLink `json:"applications,omitempty"`
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
	DeployedAt time.Time `json:"deployedAt,omitempty"`
	Notes      string    `json:"notes,omitempty"`
	State      string    `json:"state,omitempty"`
}

// ApplicationObject base resource type
type ApplicationObject struct {
	OID           string                   `json:"oid,omitempty"`
	ObjType       string                   `json:"objType,omitempty"`
	NamespaceID   string                   `json:"namespaceId,omitempty"`
	Deployment    string                   `json:"deploymentName,omitempty"`
	Server        string                   `json:"resgistryServer,omitempty"`
	ChartRegistry string                   `json:"chartRegistry,omitempty"`
	ChartName     string                   `json:"chartName,omitempty"`
	ChartVersion  string                   `json:"chartVersion,omitempty"`
	Channel       string                   `json:"channel,omitempty"`
	Username      string                   `json:"username,omitempty"`
	Password      string                   `json:"password,omitempty"`
	Config        string                   `json:"config,omitempty"`
	JSONValues    string                   `json:"jsonValues,omitempty"`
	CreatedAt     time.Time                `json:"createdAt,omitempty"`
	UpdatedAt     time.Time                `json:"updatedAt,omitempty"`
	Status        *ApplicationStatusObject `json:"status,omitempty"`
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
	OID          string    `json:"oid,omitempty"`
	ObjType      string    `json:"objType,omitempty"`
	NodePoolSize int       `json:"nodePoolSize,omitempty"`
	CreatedAt    time.Time `json:"createdAt,omitempty"`
	UpdatedAt    time.Time `json:"updatedAt,omitempty"`
	State        string    `json:"state,omitempty"`
	NamespaceID  string    `json:"namespaceId,omitempty"`
}

// DataModel the actual structure for the API's data.
type DataModel struct {
	Projects     map[string]*ProjectObject     `json:"projects,omitempty"`
	Namespaces   map[string]*NamespaceObject   `json:"namespaces,omitempty"`
	Resources    map[string]*ResourceObject    `json:"resources,omitempty"`
	Applications map[string]*ApplicationObject `json:"applications,omitempty"`
}

// Reset removes all entries from the database. Mainly intended for tests.
func (data *DataModel) Reset() {
	data.Projects = make(map[string]*ProjectObject)
	data.Namespaces = make(map[string]*NamespaceObject)
	data.Applications = make(map[string]*ApplicationObject)
	data.Resources = make(map[string]*ResourceObject)
}

// DataStore in-memory data synchronization structure for API data.
type DataStore struct {
	sync.Mutex
	archive chan bool
	data    DataModel
	persist string
}

// NewDataStore initializes a new "DataStore"
func NewDataStore(filepath string) (ds *DataStore) {
	backup, err := openDataStoreFileBackup(filepath)
	if err == nil && len(backup) > 0 {
		var dm DataModel
		dm.Reset()
		if err = json.Unmarshal(backup, &dm); err == nil {
			return &DataStore{
				archive: make(chan bool, 1),
				persist: filepath,
				data:    dm,
			}
		}
		glog.Warningf("Unmarshal of JSON initialization data from backup persistence file: %s, error: %v", filepath, err)
	}
	ds = &DataStore{archive: make(chan bool, 1), persist: filepath}
	ds.data.Reset()
	return ds
}

// Archiver - archiver's main loop
func (ds *DataStore) Archiver() {
	for {
		<-ds.archive

		ds.Lock()
		archive, err := json.Marshal(ds.data)
		ds.Unlock()
		if err != nil {
			glog.Warningf("JSON marshalling error: %v", err)
		} else {
			err := copyDataStoreFileBackup(ds.persist)
			if err != nil {
				glog.Warningf("failed to make backup copy of persistence file, error: %v", err)
				// intentinoally continue through to back up data even w/o backup.
			}
			err = ioutil.WriteFile(ds.persist, archive, 0644)
			if err != nil {
				glog.Warningf("failed to API persistence backup file, error: %v", err)
			}
		}
	}
}

// CheckedRandomHexString - generate a random hex string, and validate
// that the generated string is NOT a duplicate to a value already used.
func (ds *DataStore) CheckedRandomHexString() string {
	i := 9
	for i > 0 {
		tmp := RandHexString(OIDLength)
		if _, ok := ds.data.Projects[tmp]; !ok {
			if _, ok := ds.data.Namespaces[tmp]; !ok {
				if _, ok := ds.data.Applications[tmp]; !ok {
					if _, ok := ds.data.Resources[tmp]; !ok {
						return tmp
					}
				}
			}
		}
		i--
	}
	return ""
}

// String - strigify
func (data *DataModel) String() string {
	json, err := json.Marshal(data)
	if err != nil {
		glog.Warningf("JSON marshalling error: %v", err)
		return ""
	}
	return string(json)
}

// String - strigify
func (ds *DataStore) String() string {
	return "persistence_file: " + ds.persist + ", data_model: " + ds.data.String()
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
	ds.data.Projects[obj.OID] = &obj
	return &obj
}

// NewProject creates ProjectObject with given name and unique object id.
func (ds *DataStore) NewProject(name string) *ProjectObject {
	obj := ds.NewProjectObject()
	if obj == nil {
		return nil
	}
	obj.Name = name
	ds.archive <- true
	return obj
}

// ProjectsCollection returns all projects.
func (ds *DataStore) ProjectsCollection() []*ProjectObject {
	ds.Lock()
	defer ds.Unlock()
	i := 0
	collection := make([]*ProjectObject, len(ds.data.Projects))
	for _, proj := range ds.data.Projects {
		collection[i] = proj
		i++
	}
	return collection
}

// Project returns the rquested project object or, nil if not found.
func (ds *DataStore) Project(oid string) (*ProjectObject, bool) {
	ds.Lock()
	defer ds.Unlock()
	proj, ok := ds.data.Projects[oid]
	return proj, ok
}

// AddProject add a project object to the data store.
func (ds *DataStore) AddProject(obj ProjectObject) {
	ds.Lock()
	defer ds.Unlock()
	ds.data.Projects[obj.OID] = &obj
	ds.archive <- true
}

// DeleteProject removes the project and all subordinate objects.
func (ds *DataStore) DeleteProject(obj *ProjectObject) {
	for _, link := range ds.data.Projects[obj.OID].Namespaces {
		ds.DeleteNamespace(ds.data.Namespaces[link.OID])
	}

	ds.Lock()
	defer ds.Unlock()
	delete(ds.data.Projects, obj.OID)
	ds.archive <- true
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

	ds.data.Namespaces[obj.OID] = &obj
	return &obj
}

// NewNamespace creates NamespaceObject with given name
func (ds *DataStore) NewNamespace(name string) *NamespaceObject {
	obj := ds.NewNamespaceObject()
	if obj == nil {
		return nil
	}
	obj.Name = name
	ds.archive <- true
	return obj
}

// NamespacesCollection returns all Namespaces for a project
func (ds *DataStore) NamespacesCollection(projectOID string) []*NamespaceObject {
	proj, ok := ds.data.Projects[projectOID]
	if !ok {
		return nil
	}

	ds.Lock()
	defer ds.Unlock()
	i := 0
	collection := make([]*NamespaceObject, len(proj.Namespaces))
	for _, link := range proj.Namespaces {
		collection[i] = ds.data.Namespaces[link.OID]
		i++
	}
	return collection
}

// Namespace returns the rquested Namespace object or, nil if not found.
func (ds *DataStore) Namespace(oid string) (*NamespaceObject, bool) {
	ds.Lock()
	defer ds.Unlock()
	ns, ok := ds.data.Namespaces[oid]
	return ns, ok
}

// AddNamespace add a Namespace object to the data store.
func (ds *DataStore) AddNamespace(obj NamespaceObject) {
	ds.Lock()
	defer ds.Unlock()
	ds.data.Namespaces[obj.OID] = &obj
	ds.archive <- true
}

// DeleteNamespace removes the Namespace and all subordinate objects.
func (ds *DataStore) DeleteNamespace(obj *NamespaceObject) {
	for _, link := range ds.data.Namespaces[obj.OID].Applications {
		ds.DeleteApplication(ds.data.Applications[link.OID])
	}
	ds.DeleteResource(ds.data.Namespaces[obj.OID].Resources.OID)

	ds.Lock()
	defer ds.Unlock()
	delete(ds.data.Namespaces, obj.OID)
	ds.archive <- true
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
	ds.data.Applications[obj.OID] = &obj
	return &obj
}

// NewApplication creates a new application resource.
func (ds *DataStore) NewApplication(namespace, deployment, server, registry, name, version string, channel,
	username, password, config, jsonValues *string) *ApplicationObject {
	obj := ds.NewApplicationObject(namespace)
	if obj == nil {
		return nil
	}
	obj.Deployment = deployment
	obj.Server = server
	obj.ChartRegistry = registry
	obj.ChartName = name
	obj.ChartVersion = version
	if channel != nil {
		obj.Channel = *channel
	}
	if username != nil {
		obj.Username = *username
	}
	if password != nil {
		obj.Password = *password
	}
	// Optional fields, may be unset (nil)
	if config != nil {
		obj.Config = *config
	}
	if jsonValues != nil {
		obj.JSONValues = *jsonValues
	}
	ds.archive <- true
	return obj
}

// Application returns the app with the given oid if found
func (ds *DataStore) Application(oid string) (*ApplicationObject, bool) {
	ds.Lock()
	defer ds.Unlock()
	app, ok := ds.data.Applications[oid]
	return app, ok
}

// ApplicationsCollection return the collection of applications from the indicated namespace.
func (ds *DataStore) ApplicationsCollection(nsOID string) []*ApplicationObject {
	namespace, ok := ds.data.Namespaces[nsOID]
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
	ds.data.Applications[obj.OID] = &obj
	ds.archive <- true
}

// DeleteApplication deletes specified application
func (ds *DataStore) DeleteApplication(obj *ApplicationObject) {
	ds.Lock()
	defer ds.Unlock()
	delete(ds.data.Applications, obj.OID)
	ds.archive <- true
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
	ds.data.Resources[obj.OID] = &obj
	return &obj
}

// NewResource creates a new ResourceObject resource.
func (ds *DataStore) NewResource(namespace string, nodes int) *ResourceObject {
	obj := ds.NewResourceObject(namespace)
	if obj == nil {
		return nil
	}
	obj.NodePoolSize = nodes
	ds.archive <- true
	return obj
}

// Resource returns the app with the given oid if found
func (ds *DataStore) Resource(oid string) (*ResourceObject, bool) {
	ds.Lock()
	defer ds.Unlock()
	res, ok := ds.data.Resources[oid]
	return res, ok
}

// ResourceObject return the resource object from the indicated namespace.
func (ds *DataStore) ResourceObject(nsOID string) (*ResourceObject, bool) {
	ns, ok := ds.data.Namespaces[nsOID]
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
	ds.data.Resources[obj.OID] = obj
	ds.archive <- true
}

// DeleteResource deletes specified application
func (ds *DataStore) DeleteResource(resOID string) {
	ds.Lock()
	defer ds.Unlock()
	delete(ds.data.Resources, resOID)
	ds.archive <- true
}

// TODO - stick the next 3 functions in a utils package...
func copyFile(dst, src string, perm os.FileMode) error {
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()
	tmp, err := ioutil.TempFile(filepath.Dir(dst), "")
	if err != nil {
		return err
	}
	_, err = io.Copy(tmp, in)
	if err != nil {
		tmp.Close()
		os.Remove(tmp.Name())
		return err
	}
	if err = tmp.Close(); err != nil {
		os.Remove(tmp.Name())
		return err
	}
	if err = os.Chmod(tmp.Name(), perm); err != nil {
		os.Remove(tmp.Name())
		return err
	}
	return os.Rename(tmp.Name(), dst)
}

func copyDataStoreFileBackup(path string) error {
	src := path
	dest := path + "." + strconv.FormatInt(time.Now().Unix(), 10)
	if _, err := os.Stat(src); os.IsNotExist(err) {
		return err
	}
	if err := copyFile(dest, src, os.FileMode(0644)); err != nil {
		return err
	}
	return nil
}

func openDataStoreFileBackup(src string) ([]byte, error) {
	if _, err := os.Stat(src); os.IsNotExist(err) {
		os.OpenFile(src, os.O_RDONLY|os.O_CREATE, 0666)
		return nil, nil
	}
	return ioutil.ReadFile(src)
}
