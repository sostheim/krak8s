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
	Attribute("namespace_id", String, func() {
		Description("The related namespace's generated unique id, not the namespace's name")
		Example("da9871c7")
	})
	Required("name", "version", "namespace_id")
})
