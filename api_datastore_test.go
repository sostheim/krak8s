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
