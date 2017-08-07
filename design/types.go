package design

import (
	. "github.com/goadesign/goa/design"
	. "github.com/goadesign/goa/design/apidsl"
)

// Common Types and MediaTypes used across the API resources.

// Project is the project resource media type.
var Project = MediaType("application/project+json", func() {
	Description("Users and tennants of the system are represented as the type Project")
	Attributes(func() {
		Attribute("id", String, "identity of project", func() {
			Example("newco")
			MinLength(2)
		})
		Attribute("href", String, "API href of project", func() {
			Example("/projects/newco")
		})
		Attribute("created_at", DateTime, "Date of creation")
		Required("id", "href", "created_at")
	})

	View("default", func() {
		Attribute("id")
		Attribute("href")
		Attribute("created_at")
	})

	View("link", func() {
		Attribute("id")
		Attribute("href")
	})
})

// Namespace is the namespace resource media type.
var Namespace = MediaType("application/namespace+json", func() {
	Description("Users and tennants of the system are represented as the type Project")
	Attributes(func() {
		Attribute("name", String, "namespace name", func() {
			Example("newco-prod")
			MinLength(2)
		})
		Attribute("href", String, "API href of the namespace", func() {
			Example("/projects/newco/ns/newco-prod")
		})
		Attribute("created_at", DateTime, "Date of creation")
		Required("name", "href", "created_at")
	})

	View("default", func() {
		Attribute("name")
		Attribute("href")
		Attribute("created_at")
	})

	View("link", func() {
		Attribute("name")
		Attribute("href")
	})
})

// MongoPostBody is the HTTP request body type.
var MongoPostBody = Type("MongoPostBody", func() {
	Attribute("application", String, func() {
		Description("Appplication Registry Identifier")
		Default("quay.io/samsung_cnct/mongodb-replicaset")
	})
	Attribute("version", String, func() {
		Description("Appplication Version")
		Default("v1.2.0")
	})
	Required("application", "version")
})

// Mongo is the MongoDB resource's MediaType.
var Mongo = MediaType("application/mongo+json", func() {
	Description("MongoDB ReplicaSet instance representation type")
	Attributes(func() {
		Attribute("application", String, "Application registry identifier")
		Attribute("version", String, "Application version")
		Attribute("created_at", DateTime, "Date of creation")
		Attribute("state", func() {
			Description("Lifecycle state")
			Enum("create_requested", "starting", "active", "delete_requested", "deleting", "deleted")
		})
		Required("application", "version", "state", "created_at")
	})

	View("default", func() {
		Attribute("application")
		Attribute("version")
		Attribute("created_at")
	})
})

// ApplicationPostBody is the HTTP POST Request body type.
var ApplicationPostBody = Type("ApplicationPostBody", func() {
	Attribute("name", String, func() {
		Description("Application name")
	})
	Attribute("version", String, func() {
		Description("Application version string")
	})
	Attribute("set", String, func() {
		Description("Application config --set argument string")
	})
	Attribute("registry", String, func() {
		Description("Application's registry")
	})
	Required("name", "version")
})

// Application is the application resource's MediaType.
var Application = MediaType("application/app+json", func() {
	Description("Application representation type")
	Attributes(func() {
		Attribute("name", String, "Application name")
		Attribute("version", String, "Application version")
		Attribute("config", String, "Configuration value settings string")
		Attribute("registry", String, "Application registry identifier")
		Attribute("status", func() {
			Attribute("deployed_at", DateTime, "Last deployment time")
			Attribute("state", func() {
				Description("Deployment state")
				Enum("UNKNOWN", "DEPLOYED", "DELETED", "SUPERSEDED", "FAILED", "DELETING")
			})
			Attribute("notes", String, "Application specific notification / statuses / notes (if any)")
			Required("deployed_at", "state")
		})
		Required("name", "version", "status")
	})

	View("default", func() {
		Attribute("name")
		Attribute("version")
		Attribute("status")
	})
})

// ClusterPostBody is the HTTP Post request body type to create a cluster resource
var ClusterPostBody = Type("CluterPostBody", func() {
	Attribute("nodePoolSize", Integer, func() {
		Description("The number of real nodes in the pool")
		Minimum(3)
		Maximum(11)
		Default(3)
	})
	Required("nodePoolSize")
})

// Cluster is the cluster resource's MediaType.
var Cluster = MediaType("application/cluster+json", func() {
	Description("Cluster resource representation type")
	Attributes(func() {
		Attribute("nodePoolSize", Integer, "Requested node pool size")
		Attribute("created_at", DateTime, "Date of creation")
		Attribute("state", func() {
			Description("Lifecycle state")
			Enum("create_requested", "starting", "active", "delete_requested", "deleting", "deleted")
		})
		Required("nodePoolSize", "created_at", "state")
	})

	View("default", func() {
		Attribute("nodePoolSize")
		Attribute("created_at")
		Attribute("state")
	})
})
