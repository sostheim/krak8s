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

// ChartPostBody is the HTTP POST Request body type.
var ChartPostBody = Type("ChartPostBody", func() {
	Attribute("application", String, func() {
		Description("Chart identifier")
	})
	Attribute("version", String, func() {
		Description("Chart version string")
	})
	Required("application", "version")
})

// Chart is the Helm Chart resource's MediaType.
var Chart = MediaType("application/chart+json", func() {
	Description("Helm chart representation type")
	Attributes(func() {
		Attribute("application", String, "Application registry identifier")
		Attribute("version", String, "Application version")
		Required("application", "version")
	})

	View("default", func() {
		Attribute("application")
		Attribute("version")
	})
})
