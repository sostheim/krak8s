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
})

var _ = Resource("methods", func() {
	Action("deploy", func() {
		Routing(GET("deploy/:client/:namespace"))
		Description("deploy MongoDB for client to namespace")
		Params(func() {
			Param("client", String, "Left operand")
			Param("namespace", String, "Right operand")
		})
		Response(OK, "text/plain")
	})

})
