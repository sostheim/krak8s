package design

import (
	. "github.com/goadesign/goa/design"
	. "github.com/goadesign/goa/design/apidsl"
)

// ClusterPostBody is the HTTP Post request body type to create a cluster resource
var ClusterPostBody = Type("ClusterPostBody", func() {
	Attribute("nodePoolSize", Integer, func() {
		Description("The number of worker nodes in the projects resource pool")
		Minimum(3)
		Maximum(11)
		Default(3)
	})
	Attribute("namespace_id", String, func() {
		Description("The related namespace's generated unique id, not the namespace's name")
		Example("da9871c7")
	})
	Required("nodePoolSize", "namespace_id")
})

// ApplicationPostBody is the HTTP POST Request body type.
var ApplicationPostBody = Type("ApplicationPostBody", func() {
	Attribute("namespace_id", String, func() {
		Description("The related namespace's generated unique id, not the namespace's name")
		Example("da9871c7")
	})
	Attribute("deployment_name", String, func() {
		Description("Cluster application deployment name")
		Default("samsung-mongodb-replicaset")
		Example("samsung-mongodb-replicaset")
	})
	Attribute("server", String, func() {
		Description("Application chart registry host server")
		Default("quay.io")
		Example("quay.io")
	})
	Attribute("registry", String, func() {
		Description("Application chart's registry")
		Default("samsung_cnct")
		Example("samsung_cnct")
	})
	Attribute("name", String, func() {
		Description("Application chart name")
		Example("mongodb-replicaset")
	})
	Attribute("version", String, func() {
		Description("Application chart version string")
		Default("latest")
		Example("latest")
	})
	Attribute("channel", String, func() {
		Description("Application chart's channel")
		Example("stable")
	})
	Attribute("username", String, func() {
		Description("Registry server username")
	})
	Attribute("password", String, func() {
		Description("Registry server password")
	})
	Attribute("set", String, func() {
		Description("Application chart config --set argument string")
	})
	Attribute("json_values", String, func() {
		Description("Application chart's json values string")
	})
	Required("deployment_name", "name", "version", "namespace_id")
})
