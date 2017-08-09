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
		Attribute("id", String, "generated resource unique id (8 character hexadecimal value)", func() {
			Example("30299bea")
		})
		Attribute("type", String, "constant: object type", func() {
			Example("project")
		})
		Attribute("name", String, "name of project", func() {
			Example("newco")
			MinLength(2)
		})
		Attribute("created_at", DateTime, "Date of creation")
		Required("id", "type", "name", "created_at")

		Attribute("namespaces", func() {
			Attribute("id", String, "generated resource unique id")
			Attribute("url", String, "namespaces collection url", func() {
				Example("/v1/project/30299bea/namespaces")
			})
			Required("id", "url")
		})

		Attribute("resources", func() {
			Attribute("id", String, "generated resource unique id")
			Attribute("url", String, "resources object url", func() {
				Example("/v1/project/30299bea/resources")
			})
			Required("id", "url")
		})

		Attribute("applications", func() {
			Attribute("id", String, "generated resource unique id")
			Attribute("url", String, "applications collection url", func() {
				Example("/v1/project/30299bea/applications")
			})
			Required("id", "url")
		})
	})

	View("default", func() {
		Attribute("id")
		Attribute("type")
		Attribute("name")
		Attribute("created_at")
		Attribute("namespaces")
		Attribute("resources")
		Attribute("applications")
	})
})

// Namespace is the namespace resource media type.
var Namespace = MediaType("application/namespace+json", func() {
	Description("Users and tennants of the system are represented as the type Project")
	Attributes(func() {
		Attribute("id", String, "generated resource unique id (8 character hexadecimal value)", func() {
			Example("da9871c7")
		})
		Attribute("type", String, "constant: object type", func() {
			Example("namespace")
		})
		Attribute("name", String, "system wide unique namespace name", func() {
			Example("newco-prod")
			MinLength(2)
		})
		Attribute("created_at", DateTime, "Date of creation")
		Required("id", "type", "name", "created_at")
	})

	View("default", func() {
		Attribute("id")
		Attribute("type")
		Attribute("name")
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
var Application = MediaType("application/application+json", func() {
	Description("Application representation type")
	Attributes(func() {
		Attribute("id", String, "generated resource unique id (8 character hexadecimal value)", func() {
			Example("e1ea1660")
		})
		Attribute("type", String, "constant: object type", func() {
			Example("application")
		})
		Attribute("name", String, "Application name")
		Attribute("version", String, "Application version")
		Attribute("config", String, "Configuration value settings (set) string")
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
		Required("id", "type", "name", "version", "status")
	})

	View("default", func() {
		Attribute("id")
		Attribute("type")
		Attribute("name")
		Attribute("version")
		Attribute("status")
	})
})

// ClusterPostBody is the HTTP Post request body type to create a cluster resource
var ClusterPostBody = Type("CluterPostBody", func() {
	Attribute("nodePoolSize", Integer, func() {
		Description("The number of worker nodes in the projects resource pool")
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
		Attribute("id", String, "generated resource unique id (8 character hexadecimal value)", func() {
			Example("de2760b1")
		})
		Attribute("type", String, "constant: object type", func() {
			Example("cluster")
		})
		Attribute("nodePoolSize", Integer, "Requested node pool size")
		Attribute("created_at", DateTime, "Date of creation")
		Attribute("state", func() {
			Description("Lifecycle state")
			Enum("create_requested", "starting", "active", "delete_requested", "deleting", "deleted")
		})
		Required("id", "type", "nodePoolSize", "created_at", "state")
	})

	View("default", func() {
		Attribute("id")
		Attribute("type")
		Attribute("nodePoolSize")
		Attribute("created_at")
		Attribute("state")
	})
})
