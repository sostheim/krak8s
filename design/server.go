package design

import (
	. "github.com/goadesign/goa/design"
	. "github.com/goadesign/goa/design/apidsl"
)

var _ = API("krak8s", func() {
	Title("krak8s API Server")
	Description("API Service for Kraken and Kubernetes Commands")
	Host("localhost:8080")
	Scheme("http")
	BasePath("/v1")
})

var _ = Resource("mongodb", func() {
	BasePath("/mongo")

	Action("deploy", func() {
		Routing(POST("/"))
		Description("deploy MongoDB for client to namespace")
		Params(func() {
			Param("client", String, "client identifier")
			Param("ns", String, "namesace identifier")
		})
		Response(OK, "text/plain")
	})

})
