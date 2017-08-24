package design

import (
	. "github.com/goadesign/goa/design"
	. "github.com/goadesign/goa/design/apidsl"
)

// ClusterRef is the cluster resource reference media type.
var ClusterRef = MediaType("application/cluster.ref+json", func() {
	Description("An cluster reesources object reference by object id (oid), and url")
	Attributes(func() {
		Attribute("oid", String, "The cluster resources resource unique oid", func() {
			Example("de2760b1")
		})
		Attribute("url", String, "url of the collection that contains this object", func() {
			Example("/v1/project/30299bea/cluster")
		})
		Required("oid", "url")
	})

	View("default", func() {
		Attribute("oid")
		Attribute("url")
	})
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
		Attribute("updated_at", DateTime, "Date of last update")
		Attribute("state", func() {
			Description("Lifecycle state")
			Enum("create_requested", "starting", "active", "delete_requested", "deleting", "deleted")
		})
		Attribute("namespace_id", String, func() {
			Description("The related namespace's generated unique id, not the namespace's name")
			Example("da9871c7")
		})
		Required("id", "type", "nodePoolSize", "created_at", "updated_at", "state", "namespace_id")
	})

	View("default", func() {
		Attribute("id")
		Attribute("type")
		Attribute("nodePoolSize")
		Attribute("created_at")
		Attribute("updated_at")
		Attribute("state")
		Attribute("namespace_id")
	})
})

// ApplicationRef is the application resource reference media type.
var ApplicationRef = MediaType("application/application.ref+json", func() {
	Description("An application object reference by object id (oid), and url")
	Attributes(func() {
		Attribute("oid", String, "The application resource unique oid", func() {
			Example("e1ea1660")
		})
		Attribute("url", String, "url of the collection that contains this object", func() {
			Example("/v1/project/30299bea/applications")
		})
		Required("oid", "url")
	})

	View("default", func() {
		Attribute("oid")
		Attribute("url")
	})
})

// Application is the application resource's MediaType.
var Application = MediaType("application/application+json", func() {
	Description("Application deployment representation type")
	Attributes(func() {
		Attribute("id", String, "generated resource unique id (8 character hexadecimal value)", func() {
			Example("e1ea1660")
		})
		Attribute("type", String, "constant: object type", func() {
			Example("application")
		})
		Attribute("deployment_name", String, "Cluster application deployment name")
		Attribute("server", String, "Application chart registry host server")
		Attribute("registry", String, "Application registry identifier")
		Attribute("name", String, "Application chart name")
		Attribute("version", String, "Application chart version (tag) string")
		Attribute("channel", String, "Application chart's channel")
		Attribute("username", String, "Registry server username")
		Attribute("password", String, "Registry server password")
		Attribute("config", String, "Application chart config --set argument string")
		Attribute("json_values", String, "Application chart's json values stringr")
		Attribute("status", func() {
			Attribute("deployed_at", DateTime, "Last deployment time")
			Attribute("state", func() {
				Description("Deployment state")
				Enum("UNKNOWN", "DEPLOYED", "DELETED", "SUPERSEDED", "FAILED", "DELETING")
			})
			Attribute("notes", String, "Application specific notification / statuses / notes (if any)")
			Required("deployed_at", "state")
		})
		Attribute("namespace_id", String, func() {
			Description("The related namespace's generated unique id, not the namespace's name")
			Example("da9871c7")
		})
		Attribute("created_at", DateTime, "Date of creation")
		Attribute("updated_at", DateTime, "Date of last update")
		Required("id", "type", "namespace_id", "deployment_name", "server", "registry", "name", "version", "channel", "username", "password", "config", "json_values", "status", "created_at", "updated_at")
	})

	View("default", func() {
		Attribute("id")
		Attribute("type")
		Attribute("namespace_id")
		Attribute("deployment_name")
		Attribute("server")
		Attribute("registry")
		Attribute("name")
		Attribute("version")
		Attribute("channel")
		Attribute("username")
		Attribute("config")
		Attribute("json_values")
		Attribute("status")
		Attribute("created_at")
		Attribute("updated_at")
	})
})

// NamespaceRef is the namespace resource reference media type.
var NamespaceRef = MediaType("application/namespace.ref+json", func() {
	Description("Users and tennants of the system are represented as the type Project")
	Attributes(func() {
		Attribute("oid", String, "The namespace resource unique oid", func() {
			Example("da9871c7")
		})
		Attribute("url", String, "url of the collection that contains this object", func() {
			Example("/v1/project/30299bea/namespaces")
		})
		Required("oid", "url")
	})

	View("default", func() {
		Attribute("oid")
		Attribute("url")
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
		Attribute("url", String, "url of the collection that contains this object", func() {
			Example("/v1/project/30299bea/namespaces")
		})
		Attribute("name", String, "system wide unique namespace name", func() {
			Example("newco-prod")
			MinLength(2)
		})
		Attribute("created_at", DateTime, "Date of creation")

		Attribute("resources", ClusterRef, "cluster resource associated with namespace")
		Attribute("applications", CollectionOf(ApplicationRef), "applications associated with namespace")

		Required("id", "type", "name", "created_at", "resources", "applications")
	})

	View("default", func() {
		Attribute("id")
		Attribute("type")
		Attribute("name")
		Attribute("created_at")
		Attribute("resources")
		Attribute("applications")
	})
})

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

		Attribute("namespaces", CollectionOf(NamespaceRef), "namespace associations for this project")

		Required("id", "type", "name", "created_at", "namespaces")
	})

	View("default", func() {
		Attribute("id")
		Attribute("type")
		Attribute("name")
		Attribute("created_at")
		Attribute("namespaces")
	})
})
