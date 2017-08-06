package design

import (
	. "github.com/goadesign/goa/design"
	. "github.com/goadesign/goa/design/apidsl"
)

var _ = API("krak8s", func() {
	Title("krak8s API Server")
	Description("API Service for Kubernetes, Kraken, and Helm Commands")
	License(func() {
		Name("Apache-2.0")
		URL("https://github.com/samsung-cnct/krak8s/blob/master/LICENSE")
	})
	Host("localhost:8080")
	Scheme("http")
	Version("v1")
	BasePath("/v1")
	Consumes("application/json")
	Produces("application/json")
})

var _ = Resource("goa_swagger", func() {
	Description("Download the Swagger 2.0 (OpenAPI) Specification for this API")
	Files("/swagger", "swagger/swagger.json")
	Files("/swagger.json", "swagger/swagger.json")
	Files("/swagger.yaml", "swagger/swagger.yaml")
})

var _ = Resource("goa_openapi", func() {
	Description("Download the OpenAPI (Swagger 2.0) Specification for this API")
	Files("/openapi", "swagger/swagger.json")
	Files("/openapi.json", "swagger/swagger.json")
	Files("/openapi.yaml", "swagger/swagger.yaml")
})

var _ = Resource("goa_mongo", func() {
	Description("Manage {create, delete}, and check the status of MongoDB deployments")
	BasePath("/mongo")

	Action("create", func() {
		Routing(POST("/"))
		Description("Create a MongoDB for a user in a namespace")
		Payload(MongoPostBody)
		Response(OK, "application/json")
	})

	Action("read", func() {
		Routing(GET("/:user/:ns"))
		Description("Get the status of the specified user/namespace(ns) MongoDB Deloyment")
		Params(func() {
			Param("user", String, "user identity")
			Param("ns", String, "namespace identifier")
			Required("user", "ns")
		})
		Response(OK, "application/json")
	})

	Action("delete", func() {
		Routing(GET("/:user/:ns"))
		Description("Delete the user/namespace specified MongoDB Deloyment")
		Params(func() {
			Param("user", String, "user identity")
			Param("ns", String, "namespace identifier")
			Required("user", "ns")
		})
		Response(OK, "application/json")
	})
})

var MongoPostBody = Type("MongoPostBody", func() {
	Attribute("application", String, func() {
		Description("Appplication Registry Identifier")
		Default("quay.io/samsung_cnct/mongodb-replicaset")
	})
	Attribute("user", String, func() {
		Description("Associated user identity")
	})
	Attribute("namespace", String, func() {
		Description("Associated amespace identitfier")
	})
	Required("application", "user", "namespace")
})
