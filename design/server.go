package design

import (
	. "github.com/goadesign/goa/design"
	. "github.com/goadesign/goa/design/apidsl"
)

// Top level API attibutes and API specification endpoints.
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

	ResponseTemplate(Created, func(pattern string) {
		Description("The requested resource has been created")
		Status(201)
	})

	ResponseTemplate(Accepted, func(pattern string) {
		Description("The request has been accepted for processing")
		Status(202)
	})
})

// The next two resources endpoints are identical in the data they serve to the
// consumer but provide diversity in how that data is accessed to enable various
// OpenAPI / Swagger tools to fetch and render the API specification easily.

var _ = Resource("swagger", func() {
	Description("Download the Swagger 2.0 (OpenAPI) Specification for this API")
	Files("/swagger", "swagger/swagger.json")
	Files("/swagger.json", "swagger/swagger.json")
	Files("/swagger.yaml", "swagger/swagger.yaml")
})

var _ = Resource("openapi", func() {
	Description("Download the OpenAPI (Swagger 2.0) Specification for this API")
	Files("/openapi", "swagger/swagger.json")
	Files("/openapi.json", "swagger/swagger.json")
	Files("/openapi.yaml", "swagger/swagger.yaml")
})
