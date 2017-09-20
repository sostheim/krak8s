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
		glog.Infof("read initialization data from persistence file: %s", filepath)
		var dm DataModel
		dm.Reset()
		if err = json.Unmarshal(backup, &dm); err == nil {
			glog.Info("Successfully read/unmarshalled initialization JSON data from persistence file")
			return &DataStore{
				archive: make(chan bool, 1),
				persist: filepath,
				data:    dm,
			}
		}
		glog.Warningf("Unmarshal of JSON initialization data from backup persistence file: %s, error: %v", filepath, err)
	}
	glog.Infof("default initialization, no persistence data found from file: %s", filepath)
	ds = &DataStore{archive: make(chan bool, 1), persist: filepath}
	ds.data.Reset()
	return ds
}

// Archiver - archiver's main loop
func (ds *DataStore) Archiver() {
	for {
		if false == <-ds.archive {
			// exit signal
			return
		}
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
				glog.Warningf("failed to write API persistence backup file, error: %v", err)
			}
			glog.Info("successfully archived persistent API state udpate.")
		}
	}
}

// CheckedRandomHexString - generate a random hex string, and validate
// that the generated string is NOT a duplicate to a value already used.
func (ds *DataStore) CheckedRandomHexString() string {
	i := 9
	for i > 0 {
		tmp := RandHexString(OIDLength)
		if _, found := ds.data.Projects[tmp]; !found {
			if _, found := ds.data.Namespaces[tmp]; !found {
				if _, found := ds.data.Applications[tmp]; !found {
					if _, found := ds.data.Resources[tmp]; !found {
						return tmp
					}
				}
			}
		}
		i--
	}
	glog.Warning("failed to generate unique hex string value after 9 retries... stopping")
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
	ds.data.Projects[obj.OID] = &obj
	ds.Unlock()
	return &obj
}

// NewProject creates ProjectObject with given name and unique object id.
func (ds *DataStore) NewProject(name string) *ProjectObject {
	obj := ds.NewProjectObject()
	if obj == nil {
		return nil
	}
	obj.Name = name
	obj.UpdatedAt = time.Now()
	ds.archive <- true
	return obj
}

// ProjectsCollection returns all projects.
func (ds *DataStore) ProjectsCollection() []*ProjectObject {
	i := 0
	ds.Lock()
	collection := make([]*ProjectObject, len(ds.data.Projects))
	for _, proj := range ds.data.Projects {
		collection[i] = proj
		i++
	}
	ds.Unlock()
	return collection
}

// Project returns the rquested project object or, nil if not found.
func (ds *DataStore) Project(oid string) (*ProjectObject, bool) {
	ds.Lock()
	proj, ok := ds.data.Projects[oid]
	ds.Unlock()
	return proj, ok
}

// DeleteProject removes the project and all subordinate objects.
func (ds *DataStore) DeleteProject(obj *ProjectObject) {
	if _, ok := ds.data.Projects[obj.OID]; !ok {
		return
	}
	if ds.data.Projects[obj.OID].Namespaces != nil {
		for _, link := range ds.data.Projects[obj.OID].Namespaces {
			ds.DeleteNamespace(ds.data.Namespaces[link.OID])
		}
	}
	ds.Lock()
	delete(ds.data.Projects, obj.OID)
	ds.Unlock()
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
	ds.data.Namespaces[obj.OID] = &obj
	ds.Unlock()
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

	i := 0
	ds.Lock()
	collection := make([]*NamespaceObject, len(proj.Namespaces))
	for _, link := range proj.Namespaces {
		collection[i] = ds.data.Namespaces[link.OID]
		i++
	}
	ds.Unlock()
	return collection
}

// Namespace returns the rquested Namespace object or, nil if not found.
func (ds *DataStore) Namespace(oid string) (*NamespaceObject, bool) {
	ds.Lock()
	ns, ok := ds.data.Namespaces[oid]
	ds.Unlock()
	return ns, ok
}

// DeleteNamespace removes the Namespace and all subordinate objects.
func (ds *DataStore) DeleteNamespace(obj *NamespaceObject) {
	if _, ok := ds.data.Namespaces[obj.OID]; !ok {
		return
	}
	if ds.data.Namespaces[obj.OID].Applications != nil {
		for i, link := range ds.data.Namespaces[obj.OID].Applications {
			ds.DeleteApplication(ds.data.Applications[link.OID])
			ds.data.Namespaces[obj.OID].Applications[i] = nil
		}
	}
	if ds.data.Namespaces[obj.OID].Resources != nil {
		ds.DeleteResource(ds.data.Namespaces[obj.OID].Resources.OID)
		ds.data.Namespaces[obj.OID].Resources = nil
	}
	ds.Lock()
	delete(ds.data.Namespaces, obj.OID)
	ds.Unlock()
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
	ds.Lock()
	ds.data.Applications[obj.OID] = &obj
	ds.Unlock()
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
	obj.UpdatedAt = time.Now()
	ds.archive <- true
	return obj
}

// Application returns the app with the given oid if found
func (ds *DataStore) Application(oid string) (*ApplicationObject, bool) {
	ds.Lock()
	app, ok := ds.data.Applications[oid]
	ds.Unlock()
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

// DeleteApplication deletes specified application
func (ds *DataStore) DeleteApplication(obj *ApplicationObject) {
	if _, ok := ds.data.Applications[obj.OID]; !ok {
		return
	}
	ds.Lock()
	delete(ds.data.Applications, obj.OID)
	ds.Unlock()
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
	ds.data.Resources[obj.OID] = &obj
	ds.Unlock()
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
	res, ok := ds.data.Resources[oid]
	ds.Unlock()
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

// DeleteResource deletes specified application
func (ds *DataStore) DeleteResource(resOID string) {
	if _, ok := ds.data.Resources[resOID]; !ok {
		return
	}
	ds.Lock()
	delete(ds.data.Resources, resOID)
	ds.Unlock()
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
